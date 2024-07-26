CREATE TABLE IF NOT EXISTS dailyManagerRewards (
    date INTEGER NOT NULL UNIQUE,   -- MUST be consistent with year, month and day
    year INTEGER NOT NULL,          -- MUST be consistent with date
    month INTEGER NOT NULL,         -- MUST be consistent with date
    day INTEGER NOT NULL,           -- MUST be consistent with date
    iotxRewards TEXT NOT NULL,      -- uint256
    uniIotxRewards TEXT NOT NULL,   -- uint256
    exchangeRatio TEXT NOT NULL,    -- uint256
    PRIMARY KEY (year, month, day)
);
