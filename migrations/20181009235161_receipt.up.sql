CREATE TABLE `receipts` (
    `id` INT(11) unsigned auto_increment,
    `name` VARCHAR(255) DEFAULT NULL,
    `description` VARCHAR(255) DEFAULT NULL,
    `category` VARCHAR(255) DEFAULT NULL,
    `cooking_time` INT(11) NOT NULL,
    `user_id` INT(11) NOT NULL,
    `media_id` INT(11) unsigned DEFAULT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    CONSTRAINT `fk_users` FOREIGN KEY (`user_id`) REFERENCES users(`id`),
    CONSTRAINT `fk_media_receipts` FOREIGN KEY (`media_id`) REFERENCES media(`id`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `ingredients` (
    `id` INT(11) unsigned auto_increment,
    `name` VARCHAR(255) DEFAULT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `receipt_ingredients` (
    `id` INT(11) unsigned auto_increment,
    `quantity` VARCHAR(255) DEFAULT NULL,
    `receipt_id` INT(11) unsigned NOT NULL,
    `ingredient_id` INT(11) unsigned NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    CONSTRAINT `fk_receipts_receipt_ingredients`  FOREIGN KEY  (`receipt_id`) REFERENCES receipts(`id`),
    CONSTRAINT `fk_ingredients_receipt_ingredients` FOREIGN KEY (`ingredient_id`) REFERENCES ingredients(`id`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `receipt_directions` (
    `id` INT(11) unsigned auto_increment,
    `description` TEXT DEFAULT NULL,
    `receipt_id` INT(11) unsigned NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    CONSTRAINT `fk_receipts_receipt_directions` FOREIGN KEY (`receipt_id`) REFERENCES receipts(`id`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

