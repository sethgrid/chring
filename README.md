# chring

`chring` is an implementation of a consistent hash ring. It is currently a naive implementation that does not allow for virtual nodes. I plan on adding that soon.

### Usage

Create a new hash ring with `ring := chring.New()`.

You add nodes with `ring.Add(NodeName)`. I suggest you use an IP address for a node's name.

Now, you can get a consistent node destination when you `ring.Get(key)`, where `key` is any value that you want to route upon, such as a user's ID.

If a node goes down, you can `ring.Remove(NodeName)`.

### Data Visualization

You can visualize your hash ring and its node locations with `chring.Serve(ring, ":5000")`. Check it out live with `cd example; go run main.go` and load http://localhost:5000. Neat!

### Pending Development

- on visualization and in code for node manager, be able to get weights of nodes (know x% of keys in node N)
- make a node manager
- allow for virtal nodes
- allow placing of nodes at a given hash id so you can manually balance nodes
- provide a clear path for rebalancing when you add or remove a node by providing a list of nodes that require migrations of data

### Inspiration

See https://github.com/sent-hil/consistenthash