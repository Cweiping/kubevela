package cli

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/oam-dev/kubevela/apis/types"
	"github.com/oam-dev/kubevela/pkg/utils/system"
	cmdutil "github.com/oam-dev/kubevela/pkg/utils/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestConfigCommandPersistentPreRunE(t *testing.T) {
	io := cmdutil.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	fakeC := types.Args{}
	cmd := NewConfigListCommand(fakeC, io)
	assert.Nil(t, cmd.PersistentPreRunE(new(cobra.Command), []string{}))
	cmd = NewConfigGetCommand(fakeC, io)
	assert.Nil(t, cmd.PersistentPreRunE(new(cobra.Command), []string{}))
	cmd = NewConfigSetCommand(fakeC, io)
	assert.Nil(t, cmd.PersistentPreRunE(new(cobra.Command), []string{}))
	cmd = NewConfigDeleteCommand(fakeC, io)
	assert.Nil(t, cmd.PersistentPreRunE(new(cobra.Command), []string{}))
}

var _ = Describe("Test Config ", func() {
	It("make config crud test", func() {
		ctx := context.Background()
		b := bytes.Buffer{}
		io := cmdutil.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

		By("Set VELA_HOME to local")
		err := os.Setenv(system.VelaHomeEnv, ".test_vela")
		Expect(err).ToNot(HaveOccurred())
		home, err := system.GetVelaHomeDir()
		Expect(err).ToNot(HaveOccurred())
		Expect(err).ToNot(HaveOccurred())
		Expect(strings.HasSuffix(home, ".test_vela")).Should(Equal(true))
		defer os.RemoveAll(home)

		By("Create Default Env")
		env, err := GetEnv(ctx, k8sClient, "")
		Expect(err).ToNot(HaveOccurred())
		Expect(env.Name).Should(Equal("default"))

		By("vela config set test a=b")
		err = setConfig(ctx, k8sClient, "", []string{"test", "a=b"}, io)
		Expect(err).ToNot(HaveOccurred())

		By("vela config get test")
		b = bytes.Buffer{}
		io.Out = &b
		err = getConfig(ctx, k8sClient, "", []string{"test"}, io)
		Expect(err).ToNot(HaveOccurred())
		Expect(b.String()).Should(Equal("Data:\n  a: b\n"))

		By("vela config set test2 c=d")
		io.Out = os.Stdout
		err = setConfig(ctx, k8sClient, "", []string{"test2", "c=d"}, io)
		Expect(err).ToNot(HaveOccurred())

		By("vela config ls")
		b = bytes.Buffer{}
		io.Out = &b
		err = ListConfigs(ctx, k8sClient, "", io)
		Expect(err).ToNot(HaveOccurred())
		Expect(b.String()).Should(Equal("NAME \ntest \ntest2\n"))

		By("vela config del test")
		b = bytes.Buffer{}
		io.Out = &b
		err = deleteConfig(ctx, k8sClient, "", []string{"test"}, io)
		Expect(err).ToNot(HaveOccurred())

		By("vela config ls")
		b = bytes.Buffer{}
		io.Out = &b
		err = ListConfigs(ctx, k8sClient, "", io)
		Expect(err).ToNot(HaveOccurred())
		Expect(b.String()).Should(Equal("NAME \ntest2\n"))
	})
})
