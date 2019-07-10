CREATE TABLE snippets (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    state INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    seed_text TEXT NOT NULL,
    generated_text TEXT
)
CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci;

INSERT INTO snippets (state, seed_text, generated_text) VALUES
(
    2,
    "Before boarding your rocket to Mars, remember to pack these items",
    "Before boarding your rocket to Mars, remember to pack these items for your rocket's return journey:

Your first rocket. When you are ready to return to Earth, carry it aboard your capsule and bring it back to Earth with you as soon as possible. The first and last rocket will be your liftoff rocket."
);

INSERT INTO snippets (state, seed_text, generated_text) VALUES
(
    2,
    "The best way to make friends is to",
    "The best way to make friends is to help them succeed and teach them about leadership and self-help.' — Dr. Richard Wrangham, author of A Book of Self-Care

1. Don't ask questions

This one is just as important for adults ages 50 or higher. While that sounds self-evident, it's also true for children, too. 'Don't worry — it really doesn't matter. You've already answered this question. You should have,' says Dr. Eric Hirschfuss, author of The Way to a Good Life for Kids."
);

INSERT INTO snippets (state, seed_text, generated_text) VALUES
(
    3,
    "The best brunch restaurant is",
    "The best brunch restaurant is a restaurant that takes your brunch table into your own restaurant.'' ~Ruth Garber


Whether or not I'm aware of it - at Thanksgiving I actually went with the 'Festival of Trees' theme (as featured on the song) – something I've never been able to do without that one too many times. What can I say? The songs are all so different, almost everything you see or hear in the movie had to be created for this movie."
);
