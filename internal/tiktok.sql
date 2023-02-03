CREATE TABLE IF NOT EXISTS user (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(30) NOT NULL,
    password VARCHAR(30) NOT NULL,
    create_at DATETIME NOT NULL,
    update_at DATETIME NOT NULL,
    delete_at DATETIME NOT NULL
)ENGINE=InnoDB  DEFAULT CHARSET=UTF8MB4 AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS video (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    author_id BIGINT NOT NULL,
    play_url VARCHAR(255) NOT NULL,
    cover_url VARCHAR(255) NOT NULL,
    favorite_count BIGINT NOT NULL,
    comment_count BIGINT NOT NULL,
    title VARCHAR(64) NOT NULL,
    create_at DATETIME NOT NULL,
    update_at DATETIME NOT NULL,
    delete_at DATETIME NOT NULL,
    FOREIGN KEY (author_id) REFERENCES user(id)
)ENGINE=InnoDB  DEFAULT CHARSET=UTF8MB4 AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS comment (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    text VARCHAR(30) NOT NULL,
    commenter_id BIGINT NOT NULL,
    video_id BIGINT NOT NULL,
    create_at DATETIME NOT NULL,
    update_at DATETIME NOT NULL,
    delete_at DATETIME NOT NULL,
    FOREIGN KEY (commenter_id) REFERENCES user(id),
    FOREIGN KEY (video_id) REFERENCES video(id)
)ENGINE=InnoDB  DEFAULT CHARSET=UTF8MB4 AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS favorite(
    user_id BIGINT,
    video_id BIGINT,
    PRIMARY KEY (user_id, video_id),
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (video_id) REFERENCES video(id)
)ENGINE=InnoDB  DEFAULT CHARSET=UTF8MB4;
