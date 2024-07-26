-- ---------------------------------------------------------------------------------------------------------------------
-- Get Data
-- ---------------------------------------------------------------------------------------------------------------------

-- name: GetDailyManagerRewards :one
SELECT year, month, day, iotxRewards, uniIotxRewards, exchangeRatio
FROM dailyManagerRewards
WHERE date = ?;


-- ---------------------------------------------------------------------------------------------------------------------
-- List Data
-- ---------------------------------------------------------------------------------------------------------------------

-- name: ListDailyManagerRewardsByYear :many
SELECT date, year, month, day, iotxRewards, uniIotxRewards, exchangeRatio
FROM dailyManagerRewards
WHERE year = ?;

-- name: ListDailyManagerRewardsByMonth :many
SELECT date, year, month, day, iotxRewards, uniIotxRewards, exchangeRatio
FROM dailyManagerRewards
WHERE year = ? AND month = ?;


-- ---------------------------------------------------------------------------------------------------------------------
-- Insert Data
-- ---------------------------------------------------------------------------------------------------------------------

-- name: CreateDailyManagerRewards :exec
INSERT INTO dailyManagerRewards (
    date, year, month, day, iotxRewards, uniIotxRewards, exchangeRatio
) VALUES (
          ?, ?, ?, ?, ?, ?, ?
         )
    RETURNING *;


-- ---------------------------------------------------------------------------------------------------------------------
-- Update Data
-- ---------------------------------------------------------------------------------------------------------------------

-- name: UpdateDailyManagerRewards :exec
UPDATE dailyManagerRewards
set iotxRewards = ?,
    uniIotxRewards = ?,
    exchangeRatio = ?
WHERE date = ?;
