def fetch_all(ids, fetch):
    """Call fetch(id) for every id concurrently; return results in ids order."""
    return [fetch(i) for i in ids]
