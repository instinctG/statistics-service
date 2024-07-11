-- Create the Client table
CREATE TABLE Client (
                        client_name VARCHAR(255) NOT NULL,
                        exchange_name VARCHAR(255) NOT NULL,
                        label VARCHAR(255) NOT NULL,
                        pair VARCHAR(255) NOT NULL,
                        UNIQUE (client_name, exchange_name, label, pair)
);

-- Create the HistoryOrder table
CREATE TABLE HistoryOrder (
                              client_name VARCHAR(255) NOT NULL ,
                              exchange_name VARCHAR(255) NOT NULL,
                              label VARCHAR(255) NOT NULL,
                              pair VARCHAR(255) NOT NULL,
                              side VARCHAR(255) NOT NULL,
                              type VARCHAR(255) NOT NULL,
                              base_qty FLOAT8 NOT NULL,
                              price FLOAT8 NOT NULL,
                              algorithm_name_placed VARCHAR(255),
                              lowest_sell_prc FLOAT8,
                              highest_buy_prc FLOAT8,
                              commission_quote_qty FLOAT8,
                              time_placed TIMESTAMP NOT NULL,
                              CONSTRAINT fk_client
                                  FOREIGN KEY(client_name,exchange_name,label,pair)
                                      REFERENCES Client(client_name,exchange_name,label,pair)
);

-- Create the OrderBook table
CREATE TABLE OrderBook (
                           id SERIAL PRIMARY KEY,
                           exchange VARCHAR(255) NOT NULL,
                           pair VARCHAR(255) NOT NULL,
                           asks JSONB NOT NULL,
                           bids JSONB NOT NULL
);
