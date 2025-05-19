#lang scheme

; kevin zheng
; 300266080

; lecture de csv

(define (read-f filename) (call-with-input-file filename
      (lambda (input-port)
       (let loop ((line (read-line input-port)))
       (cond 
        ((eof-object? line) '())
        (#t (begin (cons (string-split line ",") (loop (read-line input-port))))))))))

; conversion du csv en liked/notLiked
(define (convert-rating L) (list (string->number (car L)) (string->number (cadr L)) (< 3.5 (string->number (caddr L)))))

; Permet de définir la liste Ratings

(define Ratings (map convert-rating (read-f "test.csv")))

; my functions:

; adds a new movie to the current list of users. If it is a rating for a new user, then this user is created.
(define (add-rating rating user-list)
(let (
(uid (car rating)) ; gets the first element (user id)
(movie (cadr rating)) ; gets the second element (movie)
(liked? (caddr rating))) ; get the 3rd element (#t for liked, #f for disliked)
(cond
      
((null? user-list) ; user doesnt exist, create new user
(list (list uid (if liked? (list movie) '())
                (if liked? '() (list movie)))))
      
((= uid (car (car user-list))) ; user exists
(cons
 (list uid
 (if liked? (cons movie (cadr (car user-list))) (cadr (car user-list)))
 (if liked? (caddr (car user-list)) (cons movie (caddr (car user-list)))))
(cdr user-list)))
      
(else (cons (car user-list)(add-rating rating (cdr user-list))))))) ; recursion


; adds all the ratings in the list to the current list of users. 
(define (add-ratings ratings user-list)
(if (null? ratings) user-list (add-ratings (cdr ratings) (add-rating (car ratings) user-list))))


; returns the user and its movies. 
(define (get-user id user-list)
  (cond
    ((null? user-list) #f)
    ((= id (car (car user-list))) (car user-list))
    (else (get-user id (cdr user-list)))))


; remove duplicates (modified function from lecture notes 04)
(define (unique L)
 (cond
 ((null? L) '())
 ((member (car L) (cdr L))
 (unique (cdr L)))
 (else (cons (car L)
 (unique (cdr L))))
) )


; finds union
; appends both lists and uses the remove duplicates function above
(define (get-union l1 l2) (unique (append l1 l2)))



; finds intersection 
(define (get-intersection l1 l2)
(cond ((null? l1) '())
((member (car l1) l2)
(cons (car l1) (get-intersection (cdr l1) l2)))
(else (get-intersection (cdr l1) l2))))



; computes the similarity between two users based on their liked and disliked movies.
; S(U1,U2)= (| L(U1) ∩ L(U2)| + | D(U1) ∩ D(U2)|) / | L(U1) ∪ L(U2) ∪ D(U1) ∪ D(U2)|
(define (get-similarity uid1 uid2)
(let* ((users (add-ratings Ratings '())) ; takes the data 
(u1 (get-user uid1 users))
(u2 (get-user uid2 users))
(likes1 (cadr u1))
(likes2 (cadr u2))
(dislikes1 (caddr u1))
(dislikes2 (caddr u2))
         
(num-likes (get-intersection likes1 likes2)) ; finds intersection of liked movies
(num-dislikes (get-intersection dislikes1 dislikes2)) ; finds intersection of disliked movies
(total (get-union (append likes1 likes2) (append dislikes1 dislikes2)))) ; finds union of all movies 
(/ (+ (length num-likes) (length num-dislikes)) (length total))))



; test
; (get-similarity 1 31)


















