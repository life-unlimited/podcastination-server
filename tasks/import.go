package tasks

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/hajimehoshi/go-mp3"
	"io/ioutil"
	"life-unlimited/podcastination/podcast_xml"
	"life-unlimited/podcastination/podcasts"
	"life-unlimited/podcastination/stores"
	"life-unlimited/podcastination/transfer"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ImportTaskDetailsFileName is the file name of the json file that holds the relevant task details.
const ImportTaskDetailsFileName = "task.json"
const PodcastXMLDetailsFileName = "podcast.xml"

// ImportJob is the task that is scheduled.
type ImportJob struct {
	PullDir        string
	PodcastDir     string
	ImportInterval time.Duration
	Store          ImportJobStores
}

type ImportJobStores struct {
	Podcasts stores.PodcastStore
	Owners   stores.OwnerStore
	Seasons  stores.SeasonStore
	Episodes stores.EpisodeStore
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
	// Assure that the image file is png
	img := task.ImageFileName
	if img != "" && !strings.HasSuffix(img, ".png") {
		return false, fmt.Errorf("image file format must be .png")
	}
	return true, nil
}

func (job *ImportJob) interval() time.Duration {
	return job.ImportInterval
}

func (job *ImportJob) name() string {
	return "ImportJob"
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
	importSuccess := 0
	changedPodcasts := make(map[int]struct{})
	for _, task := range tasks {
		affectedPodcast, err := job.performImportTask(task)
		if err != nil {
			log.Printf("could not perform import task for %s: %v", task.Details.Title, err)
			continue
		}
		importSuccess++
		changedPodcasts[affectedPodcast.Id] = struct{}{}
	}
	log.Printf("performed import tasks: %dx success, %dx failure", importSuccess, len(tasks)-importSuccess)
	if importSuccess == 0 {
		log.Printf("no podcasts need a podcast xml refresh.")
		return nil
	}
	// Generate new podcast xml files.
	podcastXMLGenerationFeedback := make(chan bool)
	for podcastId := range changedPodcasts {
		// Generate in parallel.
		go func(podcastId int, success chan bool) {
			if err := job.refreshPodcastXML(podcastId); err != nil {
				log.Printf("could not refresh podcast xml for podcast %d: %v", podcastId, err)
				success <- false
				return
			}
			success <- true
		}(podcastId, podcastXMLGenerationFeedback)
	}
	// Wait for completion.
	podcastXMLRefreshSuccess := 0
	for remaining := len(changedPodcasts); remaining > 0; remaining-- {
		success := <-podcastXMLGenerationFeedback
		if success {
			podcastXMLRefreshSuccess++
		}
	}
	close(podcastXMLGenerationFeedback)
	log.Printf("performed podcast xml refresh: %dx success, %dx failure", podcastXMLRefreshSuccess,
		len(changedPodcasts)-podcastXMLRefreshSuccess)
	// Done.
	return nil
}

// refreshPodcastXML refreshes the podcast xml file for the given podcast.
func (job *ImportJob) refreshPodcastXML(podcastId int) error {
	// Get whole podcast content.
	creationDetails, err := job.getPodcastAsCreationDetails(podcastId)
	if err != nil {
		return fmt.Errorf("could not get creation details: %v", err)
	}
	// Generate podcast xml.
	podcastXML, err := podcast_xml.GeneratePodcastXML(creationDetails)
	if err != nil {
		return fmt.Errorf("could not generate podcast xml: %v", err)
	}
	output, err := xml.MarshalIndent(podcastXML, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal podcast xml: %v", err)
	}
	// Write podcast xml.
	podcastXMLFilePath := filepath.Join(job.PodcastDir, transfer.GetPodcastFolderName(podcastId),
		PodcastXMLDetailsFileName)
	err = ioutil.WriteFile(podcastXMLFilePath, output, 0633)
	if err != nil {
		return fmt.Errorf("could not write podcast xml: %v", err)
	}
	return nil
}

// getPodcastAsCreationDetails retrieves all creation details needed in order to generate a PodcastXML.
func (job *ImportJob) getPodcastAsCreationDetails(podcastId int) (podcast_xml.CreationDetails, error) {
	podcast, err := job.Store.Podcasts.ById(podcastId)
	if err != nil {
		return podcast_xml.CreationDetails{}, fmt.Errorf("could not get podcast %d from db: %v", podcastId, err)
	}
	owner, err := job.Store.Owners.ById(podcast.OwnerId)
	if err != nil {
		return podcast_xml.CreationDetails{}, fmt.Errorf("could not get owner %d from db: %v",
			podcast.OwnerId, err)
	}
	seasons, err := job.Store.Seasons.ByPodcast(podcastId)
	if err != nil {
		return podcast_xml.CreationDetails{}, fmt.Errorf("could not get seasons for podcast %d from db: %v",
			podcastId, err)
	}
	episodes, err := job.Store.Episodes.ByPodcast(podcastId)
	if err != nil {
		return podcast_xml.CreationDetails{}, fmt.Errorf("could not get episodes for podcast %d from db: %v",
			podcastId, err)
	}
	return podcast_xml.CreationDetails{
		Owner:    owner,
		Podcast:  podcast,
		Seasons:  seasons,
		Episodes: episodes,
	}, nil
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
// moved to its final location. However this does not perform the podcast xml file refresh.
func (job *ImportJob) performImportTask(task ImportTask) (podcasts.Podcast, error) {
	// Check the mp3 file.
	audioLength, err := validateMP3(filepath.Join(task.BaseDir, task.Details.MP3FileName))
	if err != nil {
		return podcasts.Podcast{}, fmt.Errorf("error while validating mp3 file: %v", err)
	}
	// Check image.
	if len(task.Details.ImageFileName) != 0 {
		image, err := os.Open(filepath.Join(task.BaseDir, task.Details.ImageFileName))
		if err != nil {
			return podcasts.Podcast{}, fmt.Errorf("could not open image file: %v", err)
		}
		if err = image.Close(); err != nil {
			return podcasts.Podcast{}, fmt.Errorf("could not close image file: %v", err)
		}
	}
	// Now we can check the database.
	// Get the podcast.
	podcast, err := job.Store.Podcasts.ByKey(task.Details.PodcastKey)
	if err != nil {
		return podcasts.Podcast{}, fmt.Errorf("could not get podcast (%s): %v", task.Details.PodcastKey, err)
	}
	// Get the season.
	season, err := job.Store.Seasons.ByKey(task.Details.SeasonKey, podcast.Id)
	if err != nil {
		return podcast, fmt.Errorf("could not get season %s in podcast %d: %v", task.Details.SeasonKey, podcast.Id, err)
	}
	// Get current episodes in season in order to get the latest episode number.
	episodesInSeason, err := job.Store.Episodes.BySeason(season.Id)
	if err != nil {
		return podcast, fmt.Errorf("could not get episodes in season %d: %v", season.Id, err)
	}
	episodeNum := 0
	for _, episodeInSeason := range episodesInSeason {
		if episodeInSeason.Num > episodeNum {
			episodeNum = episodeInSeason.Num
		}
	}
	episodeNum++
	// Create new episode entry and insert into db as we need the assigned id.
	episode := podcasts.Episode{
		Title:       task.Details.Title,
		Subtitle:    task.Details.Subtitle,
		Date:        task.Details.Date,
		Author:      task.Details.Author,
		Description: task.Details.Description,
		MP3Length:   audioLength,
		SeasonId:    season.Id,
		Num:         episodeNum,
		YouTubeURL:  task.Details.YouTubeURL,
		IsAvailable: false, // This will be updated to true when all files are transferred.
	}
	// Insert into db and get the inserted episode with its assigned id.
	episode, err = job.Store.Episodes.Create(episode)
	if err != nil {
		return podcast, fmt.Errorf("could not insert episode into db: %v", err)
	}
	// Get new file locations.
	fileLocations := transfer.GetEpisodeFileLocations(episode, podcast.Id)
	episode.MP3Location = fileLocations.MP3FullPath()
	episode.ImageLocation = fileLocations.ImageFullPath()
	// Transfer the files.
	err = job.performFileTransfer(episode, task, fileLocations)
	if err != nil {
		return podcast, fmt.Errorf("could not perform file transfer: %v", err)
	}
	// Set active to true in db for episode.
	episode.IsAvailable = true
	err = job.Store.Episodes.Update(episode)
	if err != nil {
		return podcast, fmt.Errorf("could not update episode data in db: %v", err)
	}
	// Podcast xml generation is done after all import tasks have been performed.
	return podcast, nil
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

// performFileTransfer transfers all episode related files to the given destination. This also deletes the task
// file.
func (job *ImportJob) performFileTransfer(episode podcasts.Episode, task ImportTask,
	fileLocations transfer.EpisodeFileLocations) error {
	// Create target directory.
	err := os.MkdirAll(filepath.Join(job.PodcastDir, fileLocations.BaseDir), 0744) // Create with read-write read read.
	if err != nil {
		return fmt.Errorf("could not create episode directory: %v", err)
	}
	// Move the files.
	// Move the mp3.
	mp3Destination := filepath.Join(job.PodcastDir, episode.MP3Location)
	err = os.Rename(filepath.Join(task.BaseDir, task.Details.MP3FileName), mp3Destination)
	if err != nil {
		return fmt.Errorf("could not move mp3 to final destination: %v", err)
	}
	// Move the image if existing.
	imageSource := filepath.Join(task.BaseDir, task.Details.ImageFileName)
	_, err = os.Stat(imageSource)
	if err == nil {
		imageDestination := filepath.Join(job.PodcastDir, episode.ImageLocation)
		err = os.Rename(imageSource, imageDestination)
		if err != nil {
			return fmt.Errorf("could not move image to final destination: %v", err)
		}
	}
	// Delete the source folder
	err = os.RemoveAll(task.BaseDir)
	if err != nil {
		return fmt.Errorf("could not delete task file: %v", err)
	}
	return nil
}
