-- nolint:concurrentIndex
CREATE INDEX IF NOT EXISTS idx_legacy ON foo (bar);
-- nolint:concurrentIndex
DROP INDEX IF EXISTS idx_old;
CREATE INDEX IF NOT EXISTS idx_no_nolint ON foo (baz);
