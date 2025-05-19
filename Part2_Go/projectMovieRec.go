// Kevin Zheng
// 300266080

// Project CSI2120/CSI2520
// Winter 2025
// Robert Laganiere, uottawa.ca

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"
)

// movies with rating greater or equal are considered 'liked'
const iLiked float64 = 3.5

// Define the Recommendation type
type Recommendation struct {
	userID     int     // recommendation for this user
	movieID    int     // recommended movie ID
	movieTitle string  // recommended movie title
	score      float32 // probability that the user will like this movie
	nUsers     int     // number of users who likes this movie
}

// get the probability that this user will like this movie
func (r Recommendation) getProbLike() float32 {
	return r.score / (float32)(r.nUsers)
}

// Define the User type
// and its list of liked items
type User struct {
	userID   int
	liked    []int // list of movies with ratings >= iLiked
	notLiked []int // list of movies with ratings < iLiked
}

func (u User) getUser() int {
	return u.userID
}

func (u *User) setUser(id int) {
	u.userID = id
}

func (u *User) addLiked(id int) {
	u.liked = append(u.liked, id)
}

func (u *User) addNotLiked(id int) {
	u.notLiked = append(u.notLiked, id)
}

// Function to read the ratings CSV file and process each row.
// The output is a map in which user ID is used as key
func readRatingsCSV(fileName string) (map[int]*User, error) {
	// Open the CSV file.
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a CSV reader.
	reader := csv.NewReader(file)

	// Read first line and skip
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	// creates the map
	users := make(map[int]*User, 1000)

	// Read all records from the CSV.
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Iterate over each record and convert the strings into integers or float.
	for _, record := range records {
		if len(record) != 4 {
			return nil, fmt.Errorf("each line must contain exactly 4 integers, but found %d", len(record))
		}

		// Parse user ID integer
		uID, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("error converting '%s' to userID integer: %v", record[0], err)
		}

		// Parse movie ID integer
		mID, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, fmt.Errorf("error converting '%s' to movieID integer: %v", record[1], err)
		}

		// Parse rating float
		r, err := strconv.ParseFloat(record[2], 64)
		if err != nil {

			return nil, fmt.Errorf("error converting '%s' to rating: %v", record[2], err)
		}

		// checks if it is a new user
		u, ok := users[uID]
		if !ok {

			u = &User{uID, nil, nil}
			users[uID] = u
		}

		// ad movie in user list
		if r >= iLiked {

			u.addLiked(mID)

		} else {

			u.addNotLiked(mID)
		}
	}

	return users, nil
}

// Function to read the movies CSV file and process each row.
// The output is a map in which user ID is used as key
func readMoviesCSV(fileName string) (map[int]string, error) {
	// Open the CSV file.
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a CSV reader.
	reader := csv.NewReader(file)

	// Read first line and skip
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	// creates the map
	movies := make(map[int]string, 1000)

	// Read all records from the CSV.
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Iterate over each record and convert the strings into integers or float.
	for _, record := range records {
		if len(record) != 3 {
			return nil, fmt.Errorf("each line must contain exactly 3 entries, but found %d", len(record))
		}

		// Parse movie ID integer
		mID, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("error converting '%s' to movieID integer: %v", record[1], err)
		}

		// record 1 is the title
		movies[mID] = record[1]
	}

	return movies, nil
}

// checks if value is in the set
func member(value int, set []int) bool {

	for _, v := range set {
		if value == v {

			return true
		}
	}

	return false
}

// generator producing Recommendation instances from movie list
func generateMovieRec(wg *sync.WaitGroup, stop <-chan bool, userID int, titles map[int]string) <-chan Recommendation {

	outputStream := make(chan Recommendation)

	go func() {
		defer func() {
			wg.Done()
		}()
		defer close(outputStream)
		// defer fmt.Println("\nFin de generateMovieRec...")
		for k, v := range titles {
			select {
			case <-stop:
				return
			case outputStream <- Recommendation{userID, k, v, 0.0, 0}:
			}
		}
	}()

	return outputStream
}

// filter that removes the recommendations for films that the user has already seen
func notSeenByUserFilter(wg *sync.WaitGroup, stop <-chan bool, ratings map[int]*User, inputStream <-chan Recommendation) <-chan Recommendation {

	outputStream := make(chan Recommendation)

	go func() {
		defer wg.Done()
		defer close(outputStream)

		for i := range inputStream {

			// if movie is either liked or disliked by user (therefore seen before)
			if member(i.movieID, ratings[i.userID].liked) || member(i.movieID, ratings[i.userID].notLiked) {
				continue
			}

			select {
			case <-stop:
				return
			case outputStream <- i:
			}
		}
	}()

	return outputStream
}

// filter that removes films that have not been liked by at least K users
func likedByAtLeastKFilter(wg *sync.WaitGroup, stop <-chan bool, k int, ratings map[int]*User, inputStream <-chan Recommendation) <-chan Recommendation {

	outputStream := make(chan Recommendation)

	go func() {
		defer wg.Done()
		defer close(outputStream)

		for i := range inputStream {

			likeCount := 0

			for _, v := range ratings {

				if member(i.movieID, v.liked) {
					likeCount += 1
				}

				if likeCount >= k {
					break
				}
			}

			if likeCount < k {
				continue
			}

			select {
			case <-stop:
				return
			case outputStream <- i:
			}
		}
	}()

	return outputStream
}

// computes the score of each recommendation. Since this is the most timeconsuming task, two instances of this stage must run in parallel
func computeScoreStage(wg *sync.WaitGroup, stop <-chan bool, ratings map[int]*User, inputStream <-chan Recommendation) <-chan Recommendation {

	outputStream := make(chan Recommendation)

	go func() {
		defer wg.Done()
		defer close(outputStream)

		for i := range inputStream { // for each recommendation that passed trhough the filters

			var score float32 = 0
			LM := 0

			for _, v := range ratings { // for each user V, V≠U

				if v.userID == i.userID {
					continue
				}

				if member(i.movieID, v.liked) { // if V liked movie M
					intersectionLike := 0
					intersectionDislike := 0

					// get intersectionLikes and intersectionDislikes values
					for _, movie := range ratings[i.userID].liked {
						if member(movie, ratings[v.userID].liked) {
							intersectionLike += 1
						}
					}
					for _, movie := range ratings[i.userID].notLiked {
						if member(movie, ratings[v.userID].notLiked) {
							intersectionDislike += 1
						}
					}

					uniqueMovies := []int{}

					//adding elements from both lists of movies from both users
					for _, movie := range ratings[i.userID].liked {
						if !member(movie, uniqueMovies) {
							uniqueMovies = append(uniqueMovies, movie)
						}
					}
					for _, movie := range ratings[i.userID].notLiked {
						if !member(movie, uniqueMovies) {
							uniqueMovies = append(uniqueMovies, movie)
						}
					}
					for _, movie := range ratings[v.userID].liked {
						if !member(movie, uniqueMovies) {
							uniqueMovies = append(uniqueMovies, movie)
						}
					}
					for _, movie := range ratings[v.userID].notLiked {
						if !member(movie, uniqueMovies) {
							uniqueMovies = append(uniqueMovies, movie)
						}
					}

					score += (float32(intersectionLike) + float32(intersectionDislike)) / float32(len(uniqueMovies))
					LM += 1
				}
			}

			// updates recommendation
			i.score = score / float32(LM)
			i.nUsers = LM

			select {
			case <-stop:
				return
			case outputStream <- i:
			}

		}
	}()

	return outputStream
}

func main() {

	fmt.Println("Number of CPUs:", runtime.NumCPU()) // just curious

	// user to be considered
	var currentUser int
	fmt.Println("Recommendations for which user? ")
	fmt.Scanf("%d", &currentUser)

	// Call the function to read and parse the movies CSV file.
	titles, err := readMoviesCSV("movies1.csv")
	if err != nil {
		log.Fatal(err)
	}

	// Call the function to read and parse the ratings CSV file.
	ratings, err := readRatingsCSV("ratings1.csv")
	if err != nil {
		log.Fatal(err)
	}

	// synchronization
	stop := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(5)

	recommendationSlice := []Recommendation{} // slice to store the recommendations

	start := time.Now() // chrono

	// the sequence of filters
	recChannel := generateMovieRec(&wg, stop, currentUser, titles)    // module 1
	filter1 := notSeenByUserFilter(&wg, stop, ratings, recChannel)    // module 2
	filter2 := likedByAtLeastKFilter(&wg, stop, 10, ratings, filter1) // module 3

	computeChannel1 := computeScoreStage(&wg, stop, ratings, filter2) // module 4A
	computeChannel2 := computeScoreStage(&wg, stop, ratings, filter2) // module 4B
	//computeChannel3 := computeScoreStage(&wg, stop, ratings, filter2)
	//computeChannel4 := computeScoreStage(&wg, stop, ratings, filter2)

	// running the goroutines...
	for {
		select {
		case recommendation, ok := <-computeChannel1:
			if ok {
				recommendationSlice = append(recommendationSlice, recommendation)
			} else {
				computeChannel1 = nil
			}
		case recommendation, ok := <-computeChannel2:
			if ok {
				recommendationSlice = append(recommendationSlice, recommendation)
			} else {
				computeChannel2 = nil
			}
			/*case recommendation, ok := <-computeChannel3:
				if ok {
					recommendationSlice = append(recommendationSlice, recommendation)
				} else {
					computeChannel3 = nil
				}
			case recommendation, ok := <-computeChannel4:
				if ok {
					recommendationSlice = append(recommendationSlice, recommendation)
				} else {
					computeChannel4 = nil
				}*/
		}
		// breaks out of the loop when all computeChannels are terminated
		if computeChannel1 == nil && computeChannel2 == nil { //&& computeChannel3 == nil && computeChannel4 == nil {
			break
		}
	}

	// sorts the slice of recommendations by score value from biggest to smallest
	sort.Slice(recommendationSlice, func(i, j int) bool {
		return recommendationSlice[i].score > recommendationSlice[j].score
	})

	close(stop) // stop all threads
	wg.Wait()   // synchronization

	end := time.Now() // ends timer

	fmt.Println("Recommendations for user", currentUser)

	// print results for top 20 recommendatio s
	for i := 0; i < 20; i++ {
		fmt.Println(recommendationSlice[i].movieTitle, " at ", recommendationSlice[i].score, " [", recommendationSlice[i].nUsers, "]")
	}
	//fmt.Println(ratings[currentUser])

	fmt.Printf("\n\nExecution time: %s", end.Sub(start))
}
