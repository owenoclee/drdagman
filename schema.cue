import "list"

#node_id: =~"^([a-z]|[A-Z]){1,32}$"

#node: {
    address!: =~".+"
    schema: {
        input: close({})
        output: close({})
    }
}

#transition: {
    when: "always"
    schema_compatability_policy: {
        enforce_compatability: "no"
    }
}

nodes!: {
    [#node_id]: #node
}

transitions!: {
    [#node_id]: {
        [#node_id]: #transition
    }
}

meta: {
    root: _graph_roots[0]
    leaf: _graph_leaves[0]
}

_node_ids: [...#node_id] & [for id, _ in nodes {id}]

_nodes_used_in_transitions: {
    let node_list = list.FlattenN([for nId, t in transitions {[nId, [for nIdP, tP in t {nIdP}]]}], -1)
    let node_map = {for n in node_list {(n): _}} 
    let node_set = list.SortStrings([for n, _ in node_map {n}])
    node_set
}
_nodes_referenced_but_not_declared: [
    for n in _nodes_used_in_transitions
    if !list.Contains(_node_ids, n)
    {n}
] & []
_nodes_declared_but_not_referenced: [
    for n in _node_ids
    if !list.Contains(_nodes_used_in_transitions, n)
    {n}
] & []

_graph_roots: [
    for n in _node_ids
    let nodes_with_incoming_transitions = list.FlattenN([for nId, t in transitions {[for nIdP, _ in t {nIdP}]}], -1)
    if !list.Contains(nodes_with_incoming_transitions, n)
    {n}
]
_graph_must_have_one_root: _graph_roots & [_]

_graph_leaves: [
    for n in _node_ids
    let nodes_with_outgoing_transitions = list.FlattenN([for nId, t in transitions {nId}], -1)
    if !list.Contains(nodes_with_outgoing_transitions, n)
    {n}
]
_graph_must_have_one_leaf: _graph_leaves & [_]

// would be nice to ensure the graph does not cycle, but not sure it can be done in cue
