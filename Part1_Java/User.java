// Kevin Zheng 300266080
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

public class User {

    private int userID;

    private int numLikedMovies; //number of instances in likedMovies
    private int numDislikedMovies;//number of instances in dislikedMovies

    protected ArrayList<Movie> likedMovies; // storing all liked movies
    protected ArrayList<Movie> dislikedMovies;// storing all disliked movies
    

    // constructor
    public User(int userID) {

        this.userID = userID;

        numLikedMovies = 0;
        numDislikedMovies = 0;

        likedMovies = new ArrayList<>();
        dislikedMovies = new ArrayList<>();

    }

    // adds movies into likedMovies 
    public void addLikes(Movie movie) {

        likedMovies.add(movie);

        movie.incrementLikes();

        numLikedMovies ++;//increments like counter of that movie

    }

    // adds movies into dislikedMovies  
    public void addDisikes(Movie movie) {

        dislikedMovies.add(movie);
        numDislikedMovies ++;

    }

    // returns id of the user
    public int getUserID() {
        return userID;
    } 

    // returns number of liked movies
    public int getNumLikedMovies() {
        return numLikedMovies;
    }

    // returns number of disliked movies
    public int getNumDislikedMovies() {
        return numDislikedMovies;
    }

    // return if the user has viewed a movie
    public boolean hasViewedMovie(int id) {

        // checks if movie is in likedMovies
        for (Movie movie : likedMovies) {
            if (movie.getID() == id) {
                return true;
            }
        }
        // checks if movie is in dislikedMovies
        for (Movie movie : dislikedMovies) {
            if (movie.getID() == id) {
                return true;
            }
        }
        return false;
    }

    // return if the user has liked a movie
    public boolean hasLikedMovie(int id) {

        // checks if movie is in likedMovies
        for (Movie movie : likedMovies) {
            if (movie.getID() == id) {
                return true;
            }
        }
        return false;
    }

    // return if the user has disliked a movie
    public boolean hasDisikedMovie(int id) {

        // checks if movie is in likedMovies
        for (Movie movie : dislikedMovies) {
            if (movie.getID() == id) {
                return true;
            }
        }
        return false;
    }


}
