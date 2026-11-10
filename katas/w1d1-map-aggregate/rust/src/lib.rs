use std::collections::HashMap;

/// Sums the values for each key across all maps in the list.
///
/// [{"m1": 20}, {"m2": 30}, {"m1": 20}] -> {"m1": 40, "m2": 30}
pub fn aggregate(list: &[HashMap<String, i64>]) -> HashMap<String, i64> {
    todo!()
}

#[cfg(test)]
mod tests {
    use super::*;

    fn m(pairs: &[(&str, i64)]) -> HashMap<String, i64> {
        pairs.iter().map(|(k, v)| (k.to_string(), *v)).collect()
    }

    #[test]
    fn interview_example() {
        let list = [m(&[("m1", 20)]), m(&[("m2", 30)]), m(&[("m1", 20)])];
        assert_eq!(aggregate(&list), m(&[("m1", 40), ("m2", 30)]));
    }

    #[test]
    fn several_keys_per_map() {
        let list = [m(&[("a", 1), ("b", 2)]), m(&[("b", 3), ("c", 4)])];
        assert_eq!(aggregate(&list), m(&[("a", 1), ("b", 5), ("c", 4)]));
    }

    #[test]
    fn zero_and_negative_totals_are_kept() {
        let list = [m(&[("x", 5)]), m(&[("x", -5)]), m(&[("y", 0)])];
        assert_eq!(aggregate(&list), m(&[("x", 0), ("y", 0)]));
    }

    #[test]
    fn empty_list() {
        assert!(aggregate(&[]).is_empty());
    }

    #[test]
    fn input_is_still_usable_after() {
        let list = vec![m(&[("m1", 20)]), m(&[("m1", 20)])];
        let _ = aggregate(&list);
        assert_eq!(list[0]["m1"], 20);
    }
}
