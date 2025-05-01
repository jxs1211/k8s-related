-- Create database (if not exists)
CREATE DATABASE IF NOT EXISTS stats_db;
USE stats_db;

-- Create application user with proper permissions (matching your secret)
CREATE USER IF NOT EXISTS 'user'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON stats_db.* TO 'user'@'%';
FLUSH PRIVILEGES;

-- Create task_instance table (optimized for KEDA scaling)
CREATE TABLE IF NOT EXISTS task_instance (
    id INT AUTO_INCREMENT PRIMARY KEY,
    task_name VARCHAR(255) NOT NULL,
    state ENUM('running', 'queued', 'failed', 'success') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_state (state)  -- Important for KEDA query performance
) ENGINE=InnoDB;

-- Insert initial data (8 running tasks)
INSERT INTO task_instance (task_name, state) VALUES
('data_processing_job_1', 'running'),
('data_processing_job_2', 'running'),
('report_generation_1', 'running'),
('report_generation_2', 'running'),
('data_sync_1', 'running'),
('data_sync_2', 'running'),
('cleanup_job_1', 'running'),
('cleanup_job_2', 'running');

-- Verification queries
-- 1. Show current tasks
SELECT * FROM task_instance ORDER BY created_at DESC;

-- 2. Show task counts by state (matches KEDA query logic)
SELECT 
    state, 
    COUNT(*) as task_count,
    CONCAT('CEIL(COUNT/6)=', CEIL(COUNT(*) / 6)) as keda_value
FROM task_instance 
WHERE state IN ('running', 'queued')
GROUP BY state;

-- 3. The exact KEDA scaling query
SELECT CEIL(COUNT(*) / 6) FROM task_instance WHERE state='running' OR state='queued';
