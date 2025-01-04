CREATE TABLE IF NOT EXISTS post(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_profile(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS comment(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_profile(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS postReact(
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    is_liked INTEGER NOT NULL,
    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES post(id),
    FOREIGN KEY (user_id) REFERENCES user_profile(id)
);
CREATE TABLE IF NOT EXISTS commentReact(
    comment_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    is_liked INTEGER NOT NULL,
    PRIMARY KEY (comment_id, user_id),
    FOREIGN KEY (comment_id) REFERENCES comment(id),
    FOREIGN KEY (user_id) REFERENCES user_profile(id)
);
CREATE TABLE IF NOT EXISTS categories(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_name TEXT UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS post_category(
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    UNIQUE (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);
CREATE VIEW IF NOT EXISTS single_post AS
SELECT p.id AS post_id,
    p.content AS post_content,
    p.created_at AS post_date,
    p.title AS post_title,
    u.username AS post_author,
    u.created_at AS joined_at,
    u.id AS post_author_id,
    COALESCE(c.comments_count, 0) AS post_comments_count,
    COALESCE(l.likes_count, 0) AS post_likes,
    COALESCE(d.dislikes_count, 0) AS post_dislikes
FROM post p
    JOIN user_profile u ON p.user_id = u.id
    LEFT JOIN (
        SELECT post_id,
            COUNT(*) AS comments_count
        FROM comment
        GROUP BY post_id
    ) c ON p.id = c.post_id
    LEFT JOIN (
        SELECT post_id,
            COUNT(*) AS likes_count
        FROM postReact
        WHERE post_id IS NOT NULL
            AND is_liked = 1
        GROUP BY post_id
    ) l ON p.id = l.post_id
    LEFT JOIN (
        SELECT post_id,
            COUNT(*) AS dislikes_count
        FROM postReact
        WHERE post_id IS NOT NULL
            AND is_liked = 2
        GROUP BY post_id
    ) d ON p.id = d.post_id;
CREATE VIEW IF NOT EXISTS single_comment AS
SELECT c.id AS comment_id,
    c.content AS comment_content,
    c.created_at AS comment_date,
    c.post_id AS post_id,
    u.username AS comment_author,
    COALESCE(l.likes_count, 0) AS comment_likes,
    COALESCE(d.dislikes_count, 0) AS comment_dislikes
FROM comment c
    JOIN user_profile u ON c.user_id = u.id
    LEFT JOIN (
        SELECT comment_id,
            COUNT(*) AS likes_count
        FROM commentReact
        WHERE comment_id IS NOT NULL
            AND is_liked = 1
        GROUP BY comment_id
    ) l ON c.id = l.comment_id
    LEFT JOIN (
        SELECT comment_id,
            COUNT(*) AS dislikes_count
        FROM commentReact
        WHERE comment_id IS NOT NULL
            AND is_liked = 2
        GROUP BY comment_id
    ) d ON c.id = d.comment_id;
INSERT
    OR IGNORE INTO categories (category_name)
VALUES ('javascript'),
    ('golang'),
    ('rust'),
    ('programming'),
    ('tech');