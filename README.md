# chring

`chring` is an implementation of a consistent hash ring. It is currently a naive implementation that does not allow for virtual nodes. I plan on adding that soon.

### Usage

Create a new hash ring with `ring := chring.New()`.

You add nodes with `ring.Add(NodeName)`. I suggest you use an IP address for a node's name.

Now, you can get a consistent node destination when you `ring.Get(key)`, where `key` is any value that you want to route upon, such as a user's ID.

If a node goes down, you can `ring.Remove(NodeName)`.

### Pending Development

- allow for virtal nodes
- allow for a visual representation of your ring
- provide a clear path for rebalancing when you add or remove a node by providing a list of nodes that require migrations of data

### Inspiration

See https://github.com/sent-hil/consistenthash