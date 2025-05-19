// Kevin Zheng 300266080

public class Recommendation {

    private User user; //the recommendation is for this user 
    private Movie recommendedMovie; // recommended movie 
    private double score; // probability that the user will like this movie 
    private int nUsers; //  number of users who likes this movie 


    // constructor
    public Recommendation(User user, Movie recommendedMovie, double score, int nUsers) {
        this.user = user;
        this.recommendedMovie = recommendedMovie;
        this.score = score;
        this.nUsers = nUsers;
    }

    // get the score of the movie recommendation
    public double getScore() {
        return score;
    }

    public String toString() {
        return recommendedMovie.toString() + " at " + score + " [ " + nUsers + "]";
    }
}
