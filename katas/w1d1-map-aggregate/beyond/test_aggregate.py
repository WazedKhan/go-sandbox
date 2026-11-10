from aggregate import aggregate


def test_interview_example():
    assert aggregate([{"m1": 20}, {"m2": 30}, {"m1": 20}]) == {"m1": 40, "m2": 30}


def test_zero_and_negative_totals_are_kept():
    assert aggregate([{"x": 5}, {"x": -5}, {"y": 0}]) == {"x": 0, "y": 0}


def test_none_inside_list():
    assert aggregate([{"a": 1}, None, {"a": 2}]) == {"a": 3}


def test_empty():
    assert aggregate([]) == {}


def test_returns_plain_dict():
    # a defaultdict leaks its behaviour to callers: result["typo"] silently becomes 0
    assert type(aggregate([{"a": 1}])) is dict


def test_input_not_mutated():
    first = {"m1": 20}
    aggregate([first, {"m1": 20}])
    assert first == {"m1": 20}
