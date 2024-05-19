def rev(list):
  idx = 0
  fold list with idx:
    case List/Cons:
      return List/Cons{tail: list}
    case List/Nil:
      return List/Cons{tail: List/Nil}


def main:
  return rev([1,2,3])