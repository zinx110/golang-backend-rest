CREATE TABLE IF NOT EXISTS products (
    `id` INT AUTO_INCREMENT ,
    `name` VARCHAR(100) NOT NULL UNIQUE,
    `description` TEXT NOT NULL,
    `image` VARCHAR(255) NOT NULL,
    `price` DECIMAL(10, 2) NOT NULL,
    `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `quantity` INT UNSIGNED NOT NULL ,
    
    PRIMARY KEY (id),
    UNIQUE KEY (name)
);