CREATE INDEX idx_job_status_name ON job_status(name);
CREATE INDEX idx_job_status_owner ON job_status(owner);
CREATE INDEX idx_job_status_phase ON job_status(phase);
CREATE INDEX idx_job_status_repo_owner ON job_status(repo_owner);
CREATE INDEX idx_job_status_repo_repo ON job_status(repo_repo);
CREATE INDEX idx_job_status_repo_host ON job_status(repo_host);
CREATE INDEX idx_job_status_repo_ref ON job_status(repo_ref);
CREATE INDEX idx_job_status_trigger_src ON job_status(trigger_src);
CREATE INDEX idx_job_status_success ON job_status(success);
CREATE INDEX idx_job_status_created ON job_status(created);
