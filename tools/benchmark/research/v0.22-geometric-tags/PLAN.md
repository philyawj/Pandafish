# Conservative geometric tag candidate

Baseline is tagged v0.21. The [assessment plan](../v0.22-next-assessment/PLAN.md)
defines every correctness, exact-work, native decision and screen gate. The
[reassessment](../v0.22-next-assessment/SELECTION_TAGS.md) freezes this mechanism
after the failed [child-storage attempt](../v0.22-child-storage/README.md).

Only additive geometric context and its move-generation call site change.
Child allocation, existing board/attack/update primitives, root/search/TT/history
policy, ordering, evaluation and clocks remain v0.21. Preserve every correction
and failed gate. No additional strength confirmation or live games run
automatically. Review a surviving candidate before committing.
