CREATE TABLE IF NOT EXISTS dailyAssetStatistics (
    date INTEGER NOT NULL UNIQUE,           -- MUST be consistent with year, month and day
    year INTEGER NOT NULL,                  -- MUST be consistent with date
    month INTEGER NOT NULL,                 -- MUST be consistent with date
    day INTEGER NOT NULL,                   -- MUST be consistent with date

    totalPending TEXT NOT NULL,             -- uint256
    totalStaked TEXT NOT NULL,              -- uint256
    totalDebts TEXT NOT NULL,               -- uint256
    exchangeRatio TEXT NOT NULL,            -- uint256

    managerRewards TEXT NOT NULL,           -- uint256
    managerRewardsUniIOTX TEXT NOT NULL,    -- uint256
    userRewards TEXT NOT NULL,              -- uint256
    userRewardsUniIOTX TEXT NOT NULL,       -- uint256

    PRIMARY KEY (year, month, day)
);