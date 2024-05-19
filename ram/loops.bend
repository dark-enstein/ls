type Tree:
  Node { ~lft, ~rgt }
  Addendum { val }

def sum(tree):
  match tree:
    case Tree/Node:
      return tree.lft + tree.rgt
    case Tree/Addendum:
      return tree.val


def enum(tree):
  idx = 0
  fold tree with idx:
    case Tree/Node:
      return Tree/Node {
        lft: tree.lft(idx * 2 + 0),
        rgt: tree.rgt(idx * 2 + 1),
      }
    case Tree/Addendum:
      return (idx, tree.val)


def main():
  return enum(
    Tree/Node {
      lft: Tree/Node { lft: Tree/Addendum { val: 4 }, rgt: Tree/Addendum { val: 7 }},
      rgt: Tree/Node { rgt: Tree/Addendum { val: 5 }, lft: Tree/Addendum { val: 6 }}
    }
  )