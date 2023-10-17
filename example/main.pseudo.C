#include<stdio.h>

int main() {
    int a, b, c;
    
    // Defining and reading values for a and b
    printf("Enter the value of a: ");
    scanf("%d", &a);
    printf("Enter the value of b: ");
    scanf("%d", &b);
    
    // Calculating c
    c = a*a - b*b;
    
    // Writing the value of c
    printf("c = %d\n", c);
    
    return 0;
}
