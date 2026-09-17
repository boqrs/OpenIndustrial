-- =============================================================================
-- Execution Result Constraints
-- =============================================================================

-- One production execution represents exactly one physical product.
-- Therefore an execution can have at most one execution result.

CREATE UNIQUE INDEX IF NOT EXISTS
    idx_execution_results_execution_id_unique
ON execution_results(execution_id);
