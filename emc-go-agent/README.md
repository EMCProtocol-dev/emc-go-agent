# EMC Node

![](https://www.edgematrix.pro/requester/static/images/4c67f2b1e2.png)

Beyond #ICP Layer2, serving as an entry of Computing power and Web3 in AI era

## EMC Go Agent
View agnet_test.go to get help.


## Basic Usage
```code
func TestCallEdgeApi_post(t *testing.T) {
	privateKey, err := crypto.BytesToECDSAPrivateKey([]byte("03b7dfc824b0cbcfe789ec0ce4571f3460befd0490e3d0d2aad8e3c07dbcce14"))
	if err != nil {
		t.Fatalf("unable to extract private key, %v", err)
	}

	a, _ := NewDefaultAgent(
		hclog.NewNullLogger(),
		privateKey,
		rpc.NewClientWithRpcUrl(rpc.TESTNET_ID, "https://oregon.edgematrix.xyz"))

	// Sample for stable diffusion
	// Api path: /sdapi/v1/txt2img
	prompt := "white cat and dog"
	data := `{"enable_hr":false,"denoising_strength":0,"firstphase_width":0,"firstphase_height":0,"hr_scale":2,"hr_upscaler":"","hr_second_pass_steps":0,"hr_resize_x":0,"hr_resize_y":0,"prompt":"%s","styles":[""],"seed":-1,"subseed":-1,"subseed_strength":0,"seed_resize_from_h":-1,"seed_resize_from_w":-1,"sampler_name":"","batch_size":1,"n_iter":1,"steps":50,"cfg_scale":7,"width":512,"height":512,"restore_faces":false,"tiling":false,"do_not_save_samples":false,"do_not_save_grid":false,"negative_prompt":"","eta":0,"s_churn":0,"s_tmax":0,"s_tmin":0,"s_noise":1,"override_settings":{},"override_settings_restore_afterwards":true,"script_args":[],"sampler_index":"Euler","script_name":"","send_images":true,"save_images":false,"alwayson_scripts":{}}`
	info, err := a.CallEdgeApi(
		"16Uiu2HAm14xAsnJHDqnQNQ2Qqo1SapdRk9j8mBKY6mghVDP9B9u5",
		"/sdapi/v1/txt2img",
		fmt.Sprintf(data, prompt),
		METHOD_POST)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("TelegramHash: ", info.TelegramHash)
		t.Log("Response: ", info.Response)
		t.Log("Err: ", info.Err)
	}
}
```

## SDK for other languages
GitHub repository: https://github.com/EMCProtocol-dev/emc_node

##  AI Sample for Stable Diffusion nodes
Address: https://6tq33-2iaaa-aaaap-qbhpa-cai.icp0.io/

GitHub repository: https://github.com/EMCProtocol-dev/EMC-SD

## Computing Node Test Tools
Address: https://57hlm-riaaa-aaaap-qbhfa-cai.icp0.io

GitHub repository: https://github.com/EMCProtocol-dev/EMC-Requester

## Tutorials
For tutorials, check https://edgematrix.pro/start

License
------
Apache 2.0
