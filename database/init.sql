CREATE TABLE deal.`deal` (
	id INT auto_increment NOT NULL,
    taker_order_id INT NOT NULL,
    maker_order_id INT NOT NULL,
    quantity INT UNSIGNED NOT NULL,
    price FLOAT NOT NULL,
	CONSTRAINT deal_PK PRIMARY KEY (id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE deal.`order` (
	id INT auto_increment NOT NULL,
	order_type INT NOT NULL COMMENT '1: buy, 2: sell',
	quantity INT UNSIGNED NOT NULL,
	remain_quantity INT UNSIGNED NOT NULL,
	price_type INT NOT NULL COMMENT '1: market, 2: limit',
	price FLOAT NOT NULL,
	is_cancel BOOL NOT NULL DEFAULT FALSE,
	CONSTRAINT order_PK PRIMARY KEY (id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_0900_ai_ci;
