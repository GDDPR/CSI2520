:- dynamic user/3, movie/2.
% K
min_liked(10).
% R
liked_th(3.5).
% N
number_of_rec(20).

% reads the user csv
read_users(Filename) :-
    csv_read_file(Filename, Data), assert_users(Data).
assert_users([]).
assert_users([row(U,_,_,_) | Rows]) :- \+number(U),!, assert_users(Rows).
assert_users([row(U,M,Rating,_) | Rows]) :- number(U),\+user(U,_,_), liked_th(R), Rating>=R,!,assert(user(U,[M],[])), assert_users(Rows).
assert_users([row(U,M,Rating,_) | Rows]) :- number(U),\+user(U,_,_), liked_th(R), Rating<R,!,assert(user(U,[],[M])), assert_users(Rows).
assert_users([row(U,M,Rating,_) | Rows]) :- number(U), liked_th(R), Rating>=R, !, retract(user(U,Liked,NotLiked)), assert(user(U,[M|Liked],NotLiked)), assert_users(Rows).
assert_users([row(U,M,Rating,_) | Rows]) :- number(U), liked_th(R), Rating<R, !, retract(user(U,Liked,NotLiked)), assert(user(U,Liked,[M|NotLiked])), assert_users(Rows).

% reads the movie csv
read_movies(Filename) :-
    csv_read_file(Filename, Rows), assert_movies(Rows).
assert_movies([]).
assert_movies([row(M,_,_) | Rows]) :- \+number(M),!, assert_movies(Rows).
assert_movies([row(M,Title,_) | Rows]) :- number(M),!, assert(movie(M,Title)), assert_movies(Rows).

% displays the first N movies ranked by score
display_first_n(_, 0) :- !.
display_first_n([], _) :- !.
display_first_n([H|T], N) :-
    writeln(H), 
    N1 is N - 1,
    display_first_n(T, N1).


% generates the (MovieTitle,Prob) pairs for all movies in the list for the user.
% takes in the user number, and list of all movies. Returns the reocmmendations im 3rd parameter
prob_movies(User, [], []) :- !.
prob_movies(User, [Movie|Movies], Rs) :-
    seen(User, Movie),  
    prob_movies(User, Movies, Rs). % if user has seen the movie, then it will skip it
prob_movies(User, [Movie|Movies], [(Title,Prob)|Rs]) :-
    prob(User, Movie, Prob),
    movie(Movie, Title),
    prob_movies(User, Movies, Rs).


% determines if the user has seen the movie. 
% takes in user and movie, and checks if they have watched it or not
seen(User, Movie) :-
    user(User, Liked, NotLiked),
    (member(Movie, Liked) ; member(Movie, NotLiked)).


% computes the probability that the user will like the movie. If the movie has not been seen by at least K users, the probability will be 0.0.
% takes in user and a movie, and spits out either 0.0 if not enough users have seem it, or gets the probability that user will like it using the index thing
prob(User, Movie, 0.0) :- %if movie not seen by at least K
    findall(U, user(U, _, _), Users),
    liked(Movie, Users, LikedUsers),
    length(LikedUsers, N),
    min_liked(K),
    N < K, !.
prob(User, Movie, P) :- %if movie has been seen by at least K
    findall(U, user(U, _, _), Users),
    liked(Movie, Users, LikedUsers),
    length(LikedUsers, N),
    min_liked(K),
    findall(S, (member(V, LikedUsers), V \= User, similarity(User, V, S)), L), 
    sum_list(L, Score),
    P is Score / N. % gets probability P BY DIVIDE BY n USER 


% extracts the users in the list who liked the movie.
% takes in movie abd the list of users and only returns a list of those who liked it
liked(_, [], []) :- !.
liked(Movie, [User|Users], [User|Rs]) :-
    user(User, Liked, _),
    member(Movie, Liked), !
    liked(Movie, Users, Rs).
liked(Movie, [_|Users], Rs) :-
    liked(Movie, Users, Rs).


% computes the similarity between two users.
% takes in 2 users and compares them using the j index
similarity(User1, User2, S) :-
    user(User1, Liked1, Disliked1),
    user(User2, Liked2, Disliked2),
    intersection(Liked1, Liked2, LikedInt),
    intersection(Disliked1, Disliked2, DislikedInt),
    append(Liked1, Disliked1, Union1),
    append(Liked2, Disliked2, Union2),
    union(Union1, Union2, Union),
    length(LikedInt, LL), % number of liked
    length(DislikedInt, LD), % number of liked
    length(Union, LU), % union
    S is (LL + LD) / LU.


% main recommendation predicate
recommendations(User) :- setof(M,L^movie(M,L),Ms), % generate list of all movie 
    prob_movies(User,Ms,Rec), % compute probabilities for all movies 
 	sort(2,@>=,Rec,Rec_Sorted), % sort by descending probabilities
	number_of_rec(N), display_first_n(Rec_Sorted,N). % display the result

init :- read_users('ratings.csv'), read_movies('movies.csv').
test(1):- similarity(33,88,S1), 291 is truncate(S1 * 10000),similarity(44,55,S2), 138 is truncate(S2 * 10000).
test(2):- prob(44,1080,P1), 122 is truncate(P1 * 10000), prob(44,1050,P2), 0 is truncate(P2).
test(3):- liked(1080, [28, 30, 32, 40, 45, 48, 49, 50], [28, 45, 50]).
test(4):- seen(32, 1080), \+seen(44, 1080).
test(5):- prob_movies(44,[1010, 1050, 1080, 2000],Rs), length(Rs,4), display(Rs).


% Kevin Zheng
% 300266080
