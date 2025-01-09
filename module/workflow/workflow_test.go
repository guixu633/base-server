package workflow

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/guixu633/base-server/module/config"
	"github.com/stretchr/testify/assert"
)

func TestTranslate(t *testing.T) {
	w := getWorkflow(t)
	raw := `The integration of artificial intelligence and automation into the workforce has sparked intense debate about the future nature of employment. While some experts paint a dystopian picture of widespread job displacement, others argue that technological advancement will create more opportunities than it eliminates. Historical precedent suggests that major technological shifts typically generate new categories of employment, even as they render certain roles obsolete.

Consider the Industrial Revolution: while it initially disrupted traditional manufacturing jobs, it ultimately led to increased productivity, economic growth, and the emergence of entirely new industries. Similarly, the digital revolution of the late 20th century created millions of jobs that would have been unimaginable just decades earlier. The key difference today lies in the unprecedented pace of technological change and its potential to affect cognitive, not just manual, tasks.

The challenge facing society isn't necessarily mass unemployment, but rather a fundamental shift in the types of skills and competencies required in the workforce. Soft skills like critical thinking, emotional intelligence, and creative problem-solving may become increasingly valuable, as these are areas where human capabilities still significantly surpass artificial intelligence. Educational systems and workforce development programs will need to evolve rapidly to prepare people for this new reality.

Moreover, the traditional concept of a linear career path may become obsolete. Workers of the future might need to be more adaptable, continuously learning and reinventing themselves as technology evolves. The gig economy and remote work arrangements, accelerated by recent global events, could become the norm rather than the exception. This shift could offer greater flexibility and work-life balance, but also presents challenges in terms of job security and benefits.`
	result, err := w.Translate(context.Background(), raw)
	assert.NoError(t, err)
	fmt.Println(result)
}

func getWorkflow(t *testing.T) *Workflow {
	cfg, err := config.LoadConfig("../../config.toml")
	assert.NoError(t, err)
	return NewWorkflow(&cfg.Workflow, &http.Client{})
}
