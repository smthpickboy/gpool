


Simple Goroutine Pool
=====================

Forked from go-playground/pool, with following improvement:

1. Remove unlimited pool and batch which I think is unnecessary.
2. Use context for cancellation.
3. Use functions to sync wait or async report result.

This package needs more testing at the moment.


License
------
Distributed under MIT License, please see license file in code for more details.
