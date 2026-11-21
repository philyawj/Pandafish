# Root PVS scouting after the parallel rejection

The two-worker root candidate saves elapsed time but fails its sole direct
screen at 50%, with mixed external rows and a failed fresh surrogate. Preserve
it without a retry, restore tagged v0.22 and reconsider single-worker work.

The current root searches each later child with [-infinity,-incumbent], while
internal PVS already scouts later children with a null window. Test the omitted
root-level scout: first root move unchanged; later roots use
[-incumbent-1,-incumbent], followed by the original full window only on a strict
improvement. Never select a fail-high child bound as an exact tie. No reductions,
parallel workers, evaluation, root restrictions, horizons, TT policy or clocks.

Require exact scores and the same first proven best root move on focused
reference controls and all 1,145 warm roots. Node counts may change but must be
deterministic. Three alternating depth-five pairs require at least 15% median
search elapsed saving. A prospective 10% node-saving gate after each pair stops
later timing if it fails. These targets are not adjusted after seeing results.

If cost survives, run the existing four native suites. Use a new final seed
2026112114/stream114 and exclude the consumed root-parallel fresh cohort as well
as earlier controls/live boards. One standard concurrent screen can follow a
promising mechanism; any surrogate exception requires matched-depth diagnosis
and a separate explicit decision, without relabeling failed gates. No automatic
larger confirmation or publication; check in before committing a viable candidate.
