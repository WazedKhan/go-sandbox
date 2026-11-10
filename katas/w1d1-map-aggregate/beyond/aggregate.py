def aggregate(maps):
    """Sum the values for each key across all dicts in the list.

    [{"m1": 20}, {"m2": 30}, {"m1": 20}] -> {"m1": 40, "m2": 30}
    A None entry in the list counts as an empty dict.
    """
    raise NotImplementedError
