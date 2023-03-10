CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(32) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME
)ENGINE=InnoDB  DEFAULT CHARSET=UTF8MB4 AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS videos (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    author_id BIGINT NOT NULL,
    play_url VARCHAR(255) NOT NULL,
    cover_url VARCHAR(255) NOT NULL,
    favorite_count BIGINT NOT NULL,
    comment_count BIGINT NOT NULL,
    title VARCHAR(64) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME ,
    FOREIGN KEY (author_id) REFERENCES users(id)
)ENGINE=InnoDB  DEFAULT CHARSET=UTF8MB4 AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS comments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    text VARCHAR(255) NOT NULL,
    commenter_id BIGINT NOT NULL,
    video_id BIGINT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME ,
    FOREIGN KEY (commenter_id) REFERENCES users(id),
    FOREIGN KEY (video_id) REFERENCES videos(id)
)ENGINE=InnoDB  DEFAULT CHARSET=UTF8MB4 AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS favorites(
    user_id BIGINT,
    video_id BIGINT,
    PRIMARY KEY (user_id, video_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (video_id) REFERENCES videos(id)
)ENGINE=InnoDB  DEFAULT CHARSET=UTF8MB4;


ALTER table users ADD INDEX name(name);
ALTER table users ADD INDEX password(password);
ALTER table videos ADD INDEX author_id(author_id);
ALTER table videos ADD INDEX created_at(created_at);
ALTER table comments ADD INDEX video_id(video_id);
