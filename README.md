# Spheres
The framework for a highly resilient secure distributed environment

Spheres intends to become a highly resilient framework for a distributed 
application, that can basically be used as a scaffold to run any 
application on top of. 

Its base layer is a set of three instances of the Spheres program,
that implement a triple redundant set of programs. If any of the 
three crash or get killed, one of the others will restart it again, 
if within the recovery cycle. The only way to truly end the programs
is to kill them all within the cycle time of three seconds

