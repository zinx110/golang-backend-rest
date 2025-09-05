CREATE TABLE IF NOT EXISTS `order_items` (
    `id` INT AUTO_INCREMENT,
    `orderId` INT NOT NULL,
    `productId` INT NOT NULL,
    `quantity` INT UNSIGNED NOT NULL,
    `price` DECIMAL(10, 2) NOT NULL,
    
    
    PRIMARY KEY (id),
    FOREIGN KEY (orderId) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (productId) REFERENCES products(id) ON DELETE CASCADE
);