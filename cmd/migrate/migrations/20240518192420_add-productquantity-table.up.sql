CREATE TABLE IF NOT EXISTS productquantity (
  `id` INT UNSIGNED NOT NULL,
  `quantity` INT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`id`) REFERENCES products(`id`)
);