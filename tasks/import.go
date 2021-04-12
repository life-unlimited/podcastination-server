package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"life-unlimited/podcastination/podcasts"
	"life-unlimited/podcastination/stores"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ImportTaskDetailsFileName is the file name of the json file that holds the relevant task details.
const ImportTaskDetailsFileName = "task.json"

// ImportJob is the task that is scheduled.
type ImportJob struct {
	PullDir    string
	PodcastDir string
	store      struct {
		podcasts stores.PodcastStore
		seasons  stores.SeasonStore
		episodes stores.EpisodeStore
	}
}

type ImportTask struct {
	BaseDir string
	Details ImportTaskDetails
}

// ImportTaskDetails is a json structure created by the user who wants to import a podcast.
type ImportTaskDetails struct {
	// PodcastKey references the target podcast.
	PodcastKey string `json:"podcast_key"`
	// SeasonKey references the target season.
	SeasonKey string `json:"season_key"`
	// Title is the title of the episode.
	Title string `json:"title"`
	// Subtitle is the subtitle for the episode.
	Subtitle string `json:"sub_title"`
	// Date is the creation date for the episode. This is also used in order to sort tasks for applying the right episode order.
	Date time.Time `json:"date"`
	// Author is the author of episode. If none provided, 'Unknown' will be used.
	Author string `json:"author"`
	// Description is an optional description for the episode.
	Description string `json:"description"`
	// MP3FileName is the file name of the mp3 file that is going to be added.
	MP3FileName string `json:"mp3_file"`
	// ImageFileName is the file name of an optional episode image.
	ImageFileName string `json:"image_file"`
	// YouTubeURL is the optional url to an youtube video.
	YouTubeURL string `json:"yt_url"`
}

// IsValid checks if the ImportTaskDetails has all needed properties in order to perform the import.
func (task *ImportTaskDetails) IsValid() (bool, error) {
	if len(task.PodcastKey) == 0 {
		return false, fmt.Errorf("no podcast key provided")
	}
	if len(task.SeasonKey) == 0 {
		return false, fmt.Errorf("no season key provided")
	}
	if len(task.Title) == 0 {
		return false, fmt.Errorf("no title provided")
	}
	if task.Date.IsZero() {
		return false, fmt.Errorf("no date provided")
	}
	if len(task.MP3FileName) == 0 {
		return false, fmt.Errorf("no mp3 file name provided")
	}
	return true, nil
}

func (job *ImportJob) importInterval() time.Duration {
	return 15 * time.Minute
}

// run runs the import tasks (yay).
func (job *ImportJob) run() error {
	// Retrieve import tasks.
	tasks, err := getImportTasks(job.PullDir)
	if err != nil {
		return fmt.Errorf("error while retrieving import tasks: %v", err)
	}
	if len(tasks) == 0 {
		return nil
	}
	// We have tasks to do.
	success := 0
	for _, task := range tasks {
		if err := job.performImportTask(task); err != nil {
			log.Printf("could not perform import task for %s: %v", task.Details.Title, err)
			continue
		}
		success++
	}
	log.Printf("performed import tasks: %dx success, %dx failure", success, len(tasks)-success)
	return nil
}

func getImportTasks(dir string) ([]ImportTask, error) {
	// Read directories in pull folder.
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not read directories from %s: %v", dir, err)
	}
	var importTasks []ImportTask
	for _, file := range fileInfo {
		if file.IsDir() {
			baseDir := filepath.Join(dir, file.Name())
			details, err := getImportTaskDetailsFromDir(baseDir)
			if err != nil {
				log.Printf("could not get import task details from directory %s: %v", file.Name(), err)
				continue
			}
			importTasks = append(importTasks, ImportTask{
				BaseDir: baseDir,
				Details: details,
			})
		}
	}
	// Sort by date ascending.
	sort.Slice(importTasks, func(i, j int) bool {
		// TODO: Test this.
		return importTasks[i].Details.Date.After(importTasks[j].Details.Date)
	})
	return importTasks, nil
}

func getImportTaskDetailsFromDir(dir string) (ImportTaskDetails, error) {
	// Open details file.
	taskDetailsFile, err := os.Open(filepath.Join(dir, ImportTaskDetailsFileName))
	if err != nil {
		return ImportTaskDetails{}, fmt.Errorf("could not open task details file: %v", err)
	}
	// Read content.
	byteValue, err := ioutil.ReadAll(taskDetailsFile)
	if err != nil {
		_ = taskDetailsFile.Close()
		return ImportTaskDetails{}, fmt.Errorf("could not read content of task details file: %v", err)
	}
	// Parse task details.
	var details ImportTaskDetails
	err = json.Unmarshal(byteValue, &details)
	if err != nil {
		_ = taskDetailsFile.Close()
		return ImportTaskDetails{}, fmt.Errorf("could not parse task details file: %v", err)
	}
	// Check if task details are valid.
	if _, err := details.IsValid(); err != nil {
		return ImportTaskDetails{}, fmt.Errorf("invalid task details file: %v", err)
	}
	return details, nil
}

// performImportTask finally performs the given task which means that the episode is inserted into the database and
// moved to its final location.
func (job *ImportJob) performImportTask(task ImportTask) error {
	// Check the mp3 file.
	audioLength, err := validateMP3(filepath.Join(task.BaseDir, task.Details.MP3FileName))
	if err != nil {
		return fmt.Errorf("error while validating mp3 file: %v", err)
	}
	// Check image.
	if len(task.Details.ImageFileName) != 0 {
		image, err := os.Open(filepath.Join(task.BaseDir, task.Details.ImageFileName))
		if err != nil {
			return fmt.Errorf("could not open image file: %v", err)
		}
		if err = image.Close(); err != nil {
			return fmt.Errorf("could not close image file: %v", err)
		}
	}
	// Now we can check the database.
	// Get the podcast.
	podcast, err := job.store.podcasts.ByKey(task.Details.PodcastKey)
	if err != nil {
		return fmt.Errorf("could not get podcast (%s): %v", task.Details.PodcastKey, err)
	}
	// Get the season.
	season, err := job.store.seasons.ByKey(task.Details.SeasonKey)
	if err != nil {
		return fmt.Errorf("could not get season (%s): %v", task.Details.SeasonKey, err)
	}
	// Create new episode entry and insert into db as we need the assigned id.
	episode := podcasts.Episode{
		Title:         task.Details.Title,
		Subtitle:      task.Details.Subtitle,
		Date:          task.Details.Date,
		Author:        task.Details.Author,
		Description:   task.Details.Description,
		ImageLocation: filepath.Join(task.BaseDir, task.Details.ImageFileName),
		MP3Location:   filepath.Join(task.BaseDIr, task.Details.MP3FileName),
		MP3Length:     0,
		SeasonId:      0,
		Num:           0,
		YouTubeURL:    "",
		IsAvailable:   false,
	}
	// Transfer the files.
	// TODO: Assure that we get the episode id.
	// Set active to true in db for episode.
	return nil
}

// validateMP3 validates an mp3 file and returns the audio length in seconds.
func validateMP3(file string) (int, error) {
	f, err := os.Open(file)
	if err != nil {
		return -1, fmt.Errorf("could not open mp3 file: %v", err)
	}
	// Decode mp3.
	d, err := mp3.NewDecoder(f)
	if err != nil {
		_ = f.Close()
		return -1, fmt.Errorf("could not create mp3 decoder: %v", err)
	}
	const sampleSize = 4
	samples := d.Length() / sampleSize
	audioLength := int(samples) / d.SampleRate()
	err = f.Close()
	if err != nil {
		return -1, fmt.Errorf("could not close mp3 file: %v", err)
	}
	return audioLength, nil
}

type episodeFileLocations struct {
	BaseDir       string
	MP3FileName   string
	ImageFileName string
}

func getEpisodeFileLocations(episode podcasts.Episode) episodeFileLocations {
	folderName := getFolderName(episode)
	cleanTitle := filepath.Clean(episode.Title)
	return episodeFileLocations{
		BaseDir:       folderName,
		MP3FileName:   fmt.Sprintf("%s.%s", strings.Replace(cleanTitle, " ", "_", -1), cleanTitle),
		ImageFileName: "thumb",
	}
}

// getFolderName returns the folder name created from the given episode.
func getFolderName(episode podcasts.Episode) string {
	timestamp := episode.Date.Format("yyyy-MM-dd_HHmmss")
	return fmt.Sprintf("%s_%d", timestamp, episode.Id)
}
