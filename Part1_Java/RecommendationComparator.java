// Kevin Zheng 300266080

import java.util.Comparator;

public class RecommendationComparator implements Comparator<Recommendation> {
    @Override
    public int compare(Recommendation r1, Recommendation r2) {
        return Double.compare(r2.getScore(), r1.getScore()); // Descending order for max heap
    }
}
