CREATE TABLE news (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    publish_date TIMESTAMP NOT NULL,
    description TEXT NOT NULL,
    cover_image_id INTEGER,
    created_by_id INTEGER,
    CONSTRAINT fk_created_by FOREIGN KEY (created_by_id) REFERENCES users(id),
    CONSTRAINT fk_cover_image FOREIGN KEY (cover_image_id) REFERENCES images(id)
);