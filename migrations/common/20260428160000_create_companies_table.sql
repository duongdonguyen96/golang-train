-- +goose Up
CREATE TABLE IF NOT EXISTS `companies` (
  `id` VARCHAR(36) NOT NULL,
  `code` VARCHAR(64) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `legal_name` VARCHAR(255) NULL,
  `tax_code` VARCHAR(64) NULL,
  `email` VARCHAR(255) NULL,
  `phone` VARCHAR(32) NULL,

  `address_line1` VARCHAR(255) NULL,
  `address_line2` VARCHAR(255) NULL,
  `city` VARCHAR(128) NULL,
  `state` VARCHAR(128) NULL,
  `postal_code` VARCHAR(32) NULL,
  `country` VARCHAR(2) NULL,

  `status` TINYINT NOT NULL DEFAULT 1,
  `metadata` JSON NULL,

  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_companies_code` (`code`),
  UNIQUE KEY `uk_companies_tax_code` (`tax_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS `companies`;
