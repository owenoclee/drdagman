import "strings"

#nodeIdentifier: string

#node: {
    address!: strings.MinRunes(1)
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
    [#nodeIdentifier]: #node
}

transitions!: {
    [#nodeIdentifier]: {
        [#nodeIdentifier]: #transition
    }
}
