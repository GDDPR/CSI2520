// Kevin Zheng 300266080
// Project CSI2120/CSI2520
// Winter 2025
// Robert Laganiere, uottawa.ca
import java.io.*;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.NoSuchElementException;
import java.util.PriorityQueue;

// this is the (incomplete) class that will generate the recommendations for a user
public class RecommendationEngine {
	
	public ArrayList<Movie> movies; // arraylist to store all the movies
	public HashMap<Integer, User> users; // hashmap to store user along their userID for quickly find users instances
	public ArrayList<User> usersArray; // arraylist of users, with same instaces as the hashmap users for looping purposes
	public PriorityQueue<Recommendation> recommendations; //Maxheap of recommendations for a user based on score (probability they will like the movie)
	
	// constructs a recommendation engine from files
	public RecommendationEngine(int userID, String movieFile, String ratingsFile) throws IOException, NumberFormatException {
		
		movies = new ArrayList<Movie>();
		readMovies(movieFile);

		users = new HashMap<Integer, User>();
		readUsers(ratingsFile);
		usersArray = new ArrayList<>(users.values()); // creates an array list of the values in user for looping purposes

		recommendations = new PriorityQueue<>(new RecommendationComparator()); // creates a maxheap of recommendations
		algorithm(userID);
	}
	
	// Reads the Movie csv file of the MovieLens dataset
	// It populates the list of Movies
    public void readMovies(String csvFile) throws IOException, NumberFormatException {

        String line;
        String delimiter = ","; // Assuming values are separated by commas

		BufferedReader br = new BufferedReader(new FileReader(csvFile)); 
		// Read each line from the CSV file
		line = br.readLine();
		
		while ((line = br.readLine()) != null && line.length() > 0) {
			// Split the line into parts using the delimiter
			String[] parts = line.split(delimiter);
			String title;
			
			// parse the ID
			int movieID= Integer.parseInt(parts[0]);
			
			if (parts.length < 3)
				throw new NumberFormatException("Error: Invalid line structure: " + line);

			// we assume that the first part is the ID
			// and the last one are genres, the rest is the title
			title= parts[1];
			if (parts.length > 3) {
				
				for (int i=2; i<parts.length-1; i++)
					title+= parts[i];
			}
			
			movies.add(new Movie(movieID,title));
		}
    }

	// Reads the ratings csv file 
	// It populates the hashmap of Users
	public void readUsers(String csvFile) throws IOException, NumberFormatException {

		String line;
        String delimiter = ","; // Assuming values are separated by commas

		BufferedReader br = new BufferedReader(new FileReader(csvFile)); 
		// Read each line from the CSV file
		line = br.readLine();
		
		while ((line = br.readLine()) != null && line.length() > 0) {
			// Split the line into parts using the delimiter
			String[] parts = line.split(delimiter);
			
			// parse the user ID
			int userId= Integer.parseInt(parts[0]);

			// parse the movie ID
			int movieId= Integer.parseInt(parts[1]);

			// parse the rating
			double rating = Double.parseDouble(parts[2]);

			// find instance of the movie
			Movie targetMovie = null;
			
			for (Movie movie : movies) {
				if (movie.getID() == movieId) {
					targetMovie = movie;
					break;
				}
			}

			if (targetMovie == null) { // if movie ID doesn't exist
				throw new NoSuchElementException();
			}

			// create new user or find an user that's already been created
			if (users.get(userId) == null) {
				users.put(userId, new User(userId));
			}

			if (rating >= 3.5) {
				users.get(userId).addLikes(targetMovie);
			} else {
				users.get(userId).addDisikes(targetMovie);
			}
		}
	}

	// takes 2 users and finds similarity between users can be estimated using the Jaccard index
	public double computePreferences(User user1, User user2) {

		int intersectionLikes = 0;
		int intersectionDislikes = 0;

		// get intersectionLikes and intersectionDislikes values
		for (Movie movie : user1.likedMovies) {
			if (user2.likedMovies.contains(movie)) {
				intersectionLikes ++;
			}
		}
		for (Movie movie : user1.dislikedMovies) {
			if (user2.dislikedMovies.contains(movie)) {
				intersectionDislikes ++;
			}
		}

		// initiating an array list that will only contain unique elements to calculate the union of the 4 sets
		ArrayList<Movie> uniqueMovies = new ArrayList<>();

		//adding elements from both lists of movies from both users
		for (Movie movie : user1.likedMovies) {
			if (!uniqueMovies.contains(movie)) {
				uniqueMovies.add(movie);
			}
		}
		for (Movie movie : user1.dislikedMovies) {
			if (!uniqueMovies.contains(movie)) {
				uniqueMovies.add(movie);
			}
		}
		for (Movie movie : user2.likedMovies) {
			if (!uniqueMovies.contains(movie)) {
				uniqueMovies.add(movie);
			}
		}
		for (Movie movie : user2.dislikedMovies) {
			if (!uniqueMovies.contains(movie)) {
				uniqueMovies.add(movie);
			}
		}

		// returns S(U1,U2)= | L(U1) ∩ L(U2)|+ | D(U1) ∩ D(U2)| / | L(U1) ∪ L(U2) ∪D(U1) ∪ D(U2)| 
		return (double) (intersectionLikes + intersectionDislikes) / (uniqueMovies.size()); 
	}

	// main algorithm. It will store all recommendations from highest to lowest inside a maxheap
	public void algorithm(int id) {

		for (Movie movie : movies) { // for each movie M in the dataset

			if (users.get(id) == null) { // checks if userID exists
				throw new NoSuchElementException();
			}

			if (!users.get(id).hasViewedMovie(movie.getID())) { // if the movie M not been viewed by U

				if (movie.getLikes() >= 10) { // if the movie M has been liked by at least 10 users

					double score = 0; // score(U,M)
					int LM = 0;

					for (User user : usersArray) { 
						if (user != users.get(id)) { // for each user V, V≠U
						
							if (user.hasLikedMovie(movie.getID())) { // if V liked movie M

								double computedPreferences = computePreferences(users.get(id), user); // compute S(U,V)

								score += computedPreferences; // score(U,M) += S(U,V) 

								LM ++;

							}

						}
					}

					// creates a new recommendation and adds it to a maxheap
					recommendations.add(new Recommendation(users.get(id), movie, score/LM, movie.getLikes()));

				}

			}
		}
		
	}

	// return a movie at a specific index in movies
	public Movie getMovie(int index) {
		return movies.get(index);
	}
	
	// returns number of movies in movies
	public int getNumberOfMovies() {
		
		return movies.size();
	}
	
	public static void main(String[] args) {
		
		
		try {

			// running the RecommendationEngine 
			RecommendationEngine rec = new RecommendationEngine(Integer.parseInt(args[0]), args[1], args[2]);
		
			
			System.out.println("Recommendations for user #  " + args[0] + ":");

		    // just printing few movies
			// prints out top 20 recommendations
			for (int i=0; i<20; i++) {
				
				System.out.println(rec.recommendations.poll()); // pulling out recommendations from maxheap
			}
			
        } catch (Exception e) {
            System.err.println("Error reading the file: " + e.getMessage());
        }
	}
}