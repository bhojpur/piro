UPDATE job_status SET data=jsonb_set(data::jsonb, '{conditions, didExecute}', 'true') WHERE data::jsonb->'conditions'->'didExecute' IS NULL;
