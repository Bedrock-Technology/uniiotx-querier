-- ---------------------------------------------------------------------------------------------------------------------
-- Get Data
-- ---------------------------------------------------------------------------------------------------------------------

-- name: GetDailyAssetStatistics :one
SELECT date, year, month, day, totalPending, totalStaked, totalDebts, exchangeRatio, managerRewards, managerRewardsUniIOTX, userRewards, userRewardsUniIOTX
FROM dailyAssetStatistics
WHERE date = ?;


-- ---------------------------------------------------------------------------------------------------------------------
-- List Data
-- ---------------------------------------------------------------------------------------------------------------------

-- name: ListDailyAssetStatisticsByYear :many
SELECT date, year, month, day, totalPending, totalStaked, totalDebts, exchangeRatio, managerRewards, managerRewardsUniIOTX, userRewards, userRewardsUniIOTX
FROM dailyAssetStatistics
WHERE year = ?;

-- name: ListDailyAssetStatisticsByMonth :many
SELECT date, year, month, day, totalPending, totalStaked, totalDebts, exchangeRatio, managerRewards, managerRewardsUniIOTX, userRewards, userRewardsUniIOTX
FROM dailyAssetStatistics
WHERE year = ? AND month = ?;

-- ---------------------------------------------------------------------------------------------------------------------
-- Insert Data
-- ---------------------------------------------------------------------------------------------------------------------

-- name: CreateDailyAssetStatistics :exec
INSERT INTO dailyAssetStatistics (
    date, year, month, day, totalPending, totalStaked, totalDebts, exchangeRatio, managerRewards, managerRewardsUniIOTX, userRewards, userRewardsUniIOTX
) VALUES (
             ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
         )
    RETURNING *;

-- ---------------------------------------------------------------------------------------------------------------------
-- Update Data
-- ---------------------------------------------------------------------------------------------------------------------

-- name: UpdateDailyAssetStatistics :exec
UPDATE dailyAssetStatistics
set totalPending = ?,
    totalStaked = ?,
    totalDebts = ?,
    exchangeRatio = ?,
    managerRewards = ?,
    managerRewardsUniIOTX = ?,
    userRewards = ?,
    userRewardsUniIOTX = ?
WHERE date = ?;
