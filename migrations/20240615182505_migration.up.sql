-- Create the swipes table
CREATE TABLE swipes (
                        id INT AUTO_INCREMENT PRIMARY KEY,
                        user_id INT NOT NULL,
                        target_id INT NOT NULL,
                        preference ENUM('YES', 'NO') NOT NULL,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        UNIQUE KEY unique_swipe (user_id, target_id)
);

-- Create the matches table
CREATE TABLE matches (
                         id INT AUTO_INCREMENT PRIMARY KEY,
                         user1_id INT NOT NULL,
                         user2_id INT NOT NULL,
                         match_id VARCHAR(64) NOT NULL,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         UNIQUE KEY unique_match (user1_id, user2_id, match_id),
                         KEY match_id_index (match_id)
);