package beforetwo

//Container is a generic container, accepting anything.
type TContainer []interface{}

//Put adds an element to the container.
func (c *TContainer) Put(elem interface{}) {
    *c = append(*c, elem)
}
//Get gets an element from the container.
func (c *TContainer) Get() interface{} {
    elem := (*c)[0]
    *c = (*c)[1:]
    return elem
}