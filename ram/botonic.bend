object H2 { x, y }

def distance(a, b):
  # print("about to work with 'a'")
  open H2: a
  open H2: b
  dx = b.x - a.x
  dy = b.y - a.y
  return (dx * dx + dy * dy) ** 0.5

def main():
  return distance(H2{x: 1, y: 2}, H2(4, 6))