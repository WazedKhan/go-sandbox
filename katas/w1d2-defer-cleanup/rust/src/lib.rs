/// Anything that must be closed after use. Closing can fail
/// (think: a file whose final flush hits a full disk).
pub trait Resource {
    fn size(&self) -> u64;
    fn close(&mut self) -> Result<(), String>;
}

/// Opens each named resource, adds up the sizes, and closes it.
///
/// Same rules as the Go version:
///   - At most one resource may be open at any moment.
///   - Every opened resource must be cleaned up, even when something fails.
///   - If opening fails, return the error with the name in it.
///   - If closing fails, return that error too, with the name in it.
pub fn sum_sizes<R, F>(mut open: F, names: &[&str]) -> Result<u64, String>
where
    R: Resource,
    F: FnMut(&str) -> Result<R, String>,
{
    let mut opened = Vec::new();
    for name in names {
        opened.push(open(name)?);
    }
    let total = opened.iter().map(|r| r.size()).sum();
    for r in opened.iter_mut() {
        let _ = r.close();
    }
    Ok(total)
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::cell::RefCell;
    use std::collections::HashMap;
    use std::rc::Rc;

    #[derive(Default)]
    struct State {
        open_now: i32,
        max_open: i32,
    }

    struct FakeRes {
        size: u64,
        close_err: Option<String>,
        closed: bool,
        state: Rc<RefCell<State>>,
    }

    impl Resource for FakeRes {
        fn size(&self) -> u64 {
            self.size
        }
        fn close(&mut self) -> Result<(), String> {
            if !self.closed {
                self.closed = true;
                self.state.borrow_mut().open_now -= 1;
            }
            match &self.close_err {
                Some(e) => Err(e.clone()),
                None => Ok(()),
            }
        }
    }

    // Dropping without close still releases it, like the OS reclaiming a
    // file handle. But any close error is lost, because drop can't return one.
    impl Drop for FakeRes {
        fn drop(&mut self) {
            if !self.closed {
                self.closed = true;
                self.state.borrow_mut().open_now -= 1;
            }
        }
    }

    struct FakeFs {
        sizes: HashMap<&'static str, u64>,
        open_err: HashMap<&'static str, &'static str>,
        close_err: HashMap<&'static str, &'static str>,
        state: Rc<RefCell<State>>,
    }

    impl FakeFs {
        fn new(sizes: &[(&'static str, u64)]) -> Self {
            FakeFs {
                sizes: sizes.iter().cloned().collect(),
                open_err: HashMap::new(),
                close_err: HashMap::new(),
                state: Rc::default(),
            }
        }
        fn opener(&self) -> impl FnMut(&str) -> Result<FakeRes, String> + '_ {
            move |name: &str| {
                if let Some(e) = self.open_err.get(name) {
                    return Err(e.to_string());
                }
                let mut s = self.state.borrow_mut();
                s.open_now += 1;
                s.max_open = s.max_open.max(s.open_now);
                Ok(FakeRes {
                    size: self.sizes[name],
                    close_err: self.close_err.get(name).map(|e| e.to_string()),
                    closed: false,
                    state: Rc::clone(&self.state),
                })
            }
        }
    }

    #[test]
    fn total() {
        let fs = FakeFs::new(&[("a", 10), ("b", 20), ("c", 30)]);
        assert_eq!(sum_sizes(fs.opener(), &["a", "b", "c"]), Ok(60));
    }

    #[test]
    fn one_open_at_a_time() {
        let fs = FakeFs::new(&[("a", 1), ("b", 1), ("c", 1), ("d", 1)]);
        sum_sizes(fs.opener(), &["a", "b", "c", "d"]).unwrap();
        assert_eq!(fs.state.borrow().max_open, 1, "max open at once");
        assert_eq!(fs.state.borrow().open_now, 0, "still open after return");
    }

    #[test]
    fn open_error_names_the_resource_and_leaks_nothing() {
        let mut fs = FakeFs::new(&[("a", 1), ("b", 1)]);
        fs.open_err.insert("c", "permission denied");
        let err = sum_sizes(fs.opener(), &["a", "b", "c", "d"]).unwrap_err();
        assert!(err.contains("permission denied"), "got {err:?}");
        assert!(err.contains('c'), "error {err:?} should name the resource");
        assert_eq!(fs.state.borrow().open_now, 0);
    }

    #[test]
    fn close_error_reaches_the_caller() {
        let mut fs = FakeFs::new(&[("a", 1), ("b", 1), ("c", 1)]);
        fs.close_err.insert("b", "disk full on flush");
        let err = sum_sizes(fs.opener(), &["a", "b", "c"]).unwrap_err();
        assert!(err.contains("disk full"), "got {err:?}");
        assert!(err.contains('b'), "error {err:?} should name the resource");
        assert_eq!(fs.state.borrow().open_now, 0);
    }
}
