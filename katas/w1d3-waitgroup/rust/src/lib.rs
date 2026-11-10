/// Calls `fetch` for every id concurrently and returns the results
/// in the same order as `ids`. Standard library only.
pub fn fetch_all<F>(ids: &[u32], fetch: F) -> Vec<String>
where
    F: Fn(u32) -> String + Sync,
{
    ids.iter().map(|&id| fetch(id)).collect()
}

/// Starts `n` threads that each add 1 to a shared counter and returns
/// the final count. It must always return `n`.
pub fn count_parallel(n: usize) -> usize {
    todo!()
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::thread::sleep;
    use std::time::{Duration, Instant};

    #[test]
    fn fetch_all_keeps_order() {
        let fetch = |id: u32| {
            sleep(Duration::from_millis((10 - id) as u64)); // later ids finish first
            format!("user-{id}")
        };
        assert_eq!(
            fetch_all(&[5, 1, 4, 2, 3], fetch),
            ["user-5", "user-1", "user-4", "user-2", "user-3"]
        );
    }

    #[test]
    fn fetch_all_is_concurrent() {
        let ids: Vec<u32> = (0..20).collect();
        let start = Instant::now();
        let got = fetch_all(&ids, |_| {
            sleep(Duration::from_millis(100));
            "ok".to_string()
        });
        let elapsed = start.elapsed();
        assert_eq!(got.len(), 20);
        assert!(elapsed < Duration::from_millis(600), "took {elapsed:?}; not concurrent");
    }

    #[test]
    fn fetch_all_empty() {
        assert!(fetch_all(&[], |_| "x".to_string()).is_empty());
    }

    #[test]
    fn count_parallel_is_exact() {
        for n in [1, 10, 1000] {
            assert_eq!(count_parallel(n), n);
        }
    }
}
