// Kevin Zheng 300266080
// Project CSI2120/CSI2520
// Winter 2025
// Robert Laganiere, uottawa.ca

public class Movie {
	
	private int movieID;
	private String title;
	private int likes;
	
	// constructs a movie
    public Movie(int id, String title) {
	
		movieID= id;
		this.title= title;
		likes = 0;

	}

	// adds +1 to likes
	protected void incrementLikes() {
		likes++;
	}

	// get the number of users who liked this movie
	public int getLikes() {
		return likes;
	}
	
	// gets the ID
	public int getID() {

		return movieID;
	}

	// get the movie title
	public String getTitle() {
	
		return title;
	}
	
	// string representation
	public String toString() {
      
       return "["+movieID+"]: "+title;	  
	}

}