# QuadTree
Go QuadTree module

# My need
I want to subdivide space to easily detect where big amount of population are located, and where are empty spaces near those population, in order to spawn events in interested locations.

# Idea of implementation
Each node must know how many samples it has
I don't want to re-instantiate a QuadTree at each loop, the quadtree can be update in place when any point is moving at a small cost, so each point must be referenced in a parallel data structure to locate his node instantly
Each empty node will have a score which determines how good it is to spawn an event inside (near players, not too far, not too close)

# TODO
- Implement a basic quadtree
- Implement a parallel data structure to locate instantly any object in the quadtree
- Implement the score of the empty nodes
- Graphical view (just a grid to send to any UI)