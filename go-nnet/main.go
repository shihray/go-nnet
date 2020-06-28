package go_nnet

import (
	"bytes"
	"encoding/gob"
	"math"
	"math/rand"
)

/////////////// Neuron Types

// Built in Neuron Types
// Sigmoid neuron has a sigmoid activation function
type Sigmoid struct{}

func (a *Sigmoid) Activate(sum float64) float64 {
	return 1.0 / (1.0 + math.Exp(-sum))
}
func (n *Sigmoid) DActivateDSum(sum, output float64) float64 {
	return output * (1 - output)
}

// Linear neuron has a linear activation function
type Linear struct{}

func (a *Linear) Activate(sum float64) float64 {
	return sum
}
func (a *Linear) DActivateDSum(sum, output float64) float64 {
	return 1.0
}

// Tanh 

type Tanh struct{}

func (a *Tanh) Activate(sum float64) float64 {
	return 1.7159 * math.Tanh(2.0/3.0*sum)
}
func (a *Tanh) DActivateDSum(sum, output float64) float64 {
	return 2.0 / 3.0 * 1.7159 * (1.0 - math.Pow(math.Tanh(2.0/3.0*sum), 2))
}

/////////////// Loss functions
// Computes the loss function and the derivative with the loss function with respect to 
// each of the inputs.
type Losser interface {
	LossAndDLossDPred(prediction []float64, truth []float64, derivative []float64) float64
}

type SquaredDistance struct{}

func (l SquaredDistance) LossAndDLossDPred(pred, truth, deriv []float64) float64 {
	// First output is the loss, second output is the derivative of the loss with respect
	// to the prediction
	loss := 0.0
	for i := range pred {
		diff := pred[i] - truth[i]
		deriv[i] = diff
		loss += math.Pow(diff, 2)
	}
	loss /= 2
	loss /= float64(len(pred))
	for i := range deriv {
		deriv[i] /= float64(len(pred))
	}
	return loss
}

type ManhattanDistance struct{}

func (m ManhattanDistance) LossAndDLossDPred(pred, truth, deriv []float64) float64 {
	loss := 0.0
	for i := range pred {
		loss += math.Abs(pred[i] - truth[i])
		if pred[i] > truth[i] {
			deriv[i] = 1.0 / float64(len(pred))
		} else if pred[i] < truth[i] {
			deriv[i] = -1.0 / float64(len(pred))
		} else {
			deriv[i] = 0
		}
	}
	loss /= float64(len(pred))
	return loss
}

////////////////////////////// Neurons

// Activator is an interface for the activation function of the neuron.
// Hopefully in the future this will allow for customization of nets to
// have neurons with custom activation functions. The activator has two
// methods, Activate, which is the actual activation function, and
// DActivateDSum which is the derivative of the activation function
// with respect to the sum. DActivateDSum takes in both the weighted
// sum and the output of Activate to save computation time if possible.
type Activator interface {
	Activate(float64) float64
	DActivateDSum(float64, float64) float64
}

// Neurons are the basic element of the neural net. They take
// in a set of inputs, compute a weighted sum of those inputs
// as set by neuron.weights, and then transforms that weighted
// sum into an alternate float64 as defined by the activation function
// The final weight is a bias term which is added at the end, so there
// should be one more weight than the number of inputs

// Need to rethink public/private of this.
type Neuron struct {
	Weights  []float64
	nWeights int
	Activator
}

func (n *Neuron) GobEncode() (buf []byte, err error) {
	w := bytes.NewBuffer(buf)
	encoder := gob.NewEncoder(w)
	err = encoder.Encode(n.Weights)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), err
}

func (n *Neuron) GobDecode(buf []byte) (err error) {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err = decoder.Decode(&n.Weights)
	if err != nil {
		return err
	}
	return err
}

// Compute the weighted sum and the activation function
func (n *Neuron) Process(input []float64) (output, sum float64) {
	for i, val := range input {
		sum += val * n.Weights[i]
	}
	sum += n.Weights[len(n.Weights)-1] //Bias term
	return n.Activate(sum), sum
}

// Construct a new neuron. Inputs the number of
// inputs to that neuron, and the specific kind of 
// reshaper that the net is.

// May need to rethink this. Getters and setters?
//func New(nInputs int, r Activater) *Neuron{
//    n := new(Neuron)
//}

func (n *Neuron) Initialize(nInputs int, r Activator) {
	n.Activator = r
	n.nWeights = nInputs + 1 // Plus one is for the bias term
	n.Weights = make([]float64, n.nWeights)
	// I'm not sure if this should be here or not
	n.InitializeWeights()
}

func (n *Neuron) InitializeWeights() {
	for i := range n.Weights {
		//n.Weights[i] = (rand.Float64() - 0.5)/(3*float64(n.nWeights))
		//n.Weights[i] = (rand.Float64() - 0.5)/(3.0*float64(n.nWeights))
		n.Weights[i] = rand.NormFloat64() * math.Pow(float64(n.nWeights), -0.5)
	}
}
